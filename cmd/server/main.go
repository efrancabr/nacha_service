package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/nacha-service/api/proto"
	"github.com/nacha-service/internal/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
)

const (
	defaultPort = ":50051"
	maxRetries  = 3
)

func main() {
	// Configure logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Starting NACHA gRPC server on port %s", defaultPort)

	// Create listener on TCP port
	lis, err := net.Listen("tcp", defaultPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Configure server options
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(unaryInterceptor),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
			Time:              1 * time.Minute,
			Timeout:           20 * time.Second,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             1 * time.Minute,
			PermitWithoutStream: true,
		}),
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer(opts...)

	// Create and register services
	nachaService := services.NewNachaService()
	pb.RegisterNachaServiceServer(grpcServer, nachaService)

	// Register health service
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	// Handle graceful shutdown
	go handleShutdown(grpcServer, healthServer)

	// Start server
	log.Printf("NACHA gRPC server is running on %s", defaultPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func handleShutdown(grpcServer *grpc.Server, healthServer *health.Server) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	sig := <-sigChan
	log.Printf("Received shutdown signal: %v", sig)

	// Set health status to not serving
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)

	// Give time for health status to propagate
	time.Sleep(time.Second)

	// Stop the server
	grpcServer.GracefulStop()
	log.Println("Server stopped gracefully")
}

func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	method := info.FullMethod

	// Add request logging
	log.Printf("Request: Method=%s", method)

	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var resp interface{}
	var err error
	var attempt int

	// Retry logic for internal errors
	for attempt = 1; attempt <= maxRetries; attempt++ {
		// Handle panic
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Panic in %s (attempt %d): %v", method, attempt, r)
					err = status.Error(codes.Internal, "Internal server error")
				}
			}()

			resp, err = handler(ctx, req)
		}()

		// Break if successful or if error is not retryable
		if err == nil || !isRetryableError(err) {
			break
		}

		if attempt < maxRetries {
			backoff := time.Duration(attempt*100) * time.Millisecond
			log.Printf("Retrying %s after %v (attempt %d)", method, backoff, attempt)
			time.Sleep(backoff)
		}
	}

	// Log completion
	duration := time.Since(start)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			log.Printf("Error: Method=%s Code=%v Message=%v Duration=%v Attempts=%d",
				method, st.Code(), st.Message(), duration, attempt)
		} else {
			log.Printf("Error: Method=%s Message=%v Duration=%v Attempts=%d",
				method, err, duration, attempt)
		}
	} else {
		log.Printf("Success: Method=%s Duration=%v Attempts=%d",
			method, duration, attempt)
	}

	return resp, err
}

func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	st, ok := status.FromError(err)
	if !ok {
		return false
	}

	switch st.Code() {
	case codes.Internal,
		codes.Unavailable,
		codes.DeadlineExceeded:
		return true
	default:
		return false
	}
}
