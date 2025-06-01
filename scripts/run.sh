#!/bin/bash

# Build the server and client
echo "Building server and client..."
go build -o bin/server cmd/server/main.go
go build -o bin/client cmd/client/main.go

# Start the server in background
echo "Starting NACHA gRPC server..."
./bin/server &
SERVER_PID=$!

# Wait for server to start
sleep 2

# Run the client
echo "Running NACHA client test..."
./bin/client

# Cleanup
echo "Cleaning up..."
kill $SERVER_PID
rm -f bin/server bin/client 