// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: api/proto/nacha.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	NachaService_ValidateFile_FullMethodName = "/nacha.NachaService/ValidateFile"
	NachaService_CreateFile_FullMethodName   = "/nacha.NachaService/CreateFile"
	NachaService_ExportFile_FullMethodName   = "/nacha.NachaService/ExportFile"
	NachaService_ViewFile_FullMethodName     = "/nacha.NachaService/ViewFile"
	NachaService_ViewDetails_FullMethodName  = "/nacha.NachaService/ViewDetails"
)

// NachaServiceClient is the client API for NachaService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NachaServiceClient interface {
	// Validate a NACHA file
	ValidateFile(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (*ValidationResponse, error)
	// Create a new NACHA file
	CreateFile(ctx context.Context, in *NachaFileRequest, opts ...grpc.CallOption) (*FileResponse, error)
	// Export NACHA file to different formats
	ExportFile(ctx context.Context, in *ExportRequest, opts ...grpc.CallOption) (*ExportResponse, error)
	// View complete file details
	ViewFile(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (*FileDetailsResponse, error)
	// View specific batch or entry details
	ViewDetails(ctx context.Context, in *DetailRequest, opts ...grpc.CallOption) (*DetailResponse, error)
}

type nachaServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNachaServiceClient(cc grpc.ClientConnInterface) NachaServiceClient {
	return &nachaServiceClient{cc}
}

func (c *nachaServiceClient) ValidateFile(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (*ValidationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ValidationResponse)
	err := c.cc.Invoke(ctx, NachaService_ValidateFile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nachaServiceClient) CreateFile(ctx context.Context, in *NachaFileRequest, opts ...grpc.CallOption) (*FileResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(FileResponse)
	err := c.cc.Invoke(ctx, NachaService_CreateFile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nachaServiceClient) ExportFile(ctx context.Context, in *ExportRequest, opts ...grpc.CallOption) (*ExportResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ExportResponse)
	err := c.cc.Invoke(ctx, NachaService_ExportFile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nachaServiceClient) ViewFile(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (*FileDetailsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(FileDetailsResponse)
	err := c.cc.Invoke(ctx, NachaService_ViewFile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nachaServiceClient) ViewDetails(ctx context.Context, in *DetailRequest, opts ...grpc.CallOption) (*DetailResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DetailResponse)
	err := c.cc.Invoke(ctx, NachaService_ViewDetails_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NachaServiceServer is the server API for NachaService service.
// All implementations must embed UnimplementedNachaServiceServer
// for forward compatibility.
type NachaServiceServer interface {
	// Validate a NACHA file
	ValidateFile(context.Context, *FileRequest) (*ValidationResponse, error)
	// Create a new NACHA file
	CreateFile(context.Context, *NachaFileRequest) (*FileResponse, error)
	// Export NACHA file to different formats
	ExportFile(context.Context, *ExportRequest) (*ExportResponse, error)
	// View complete file details
	ViewFile(context.Context, *FileRequest) (*FileDetailsResponse, error)
	// View specific batch or entry details
	ViewDetails(context.Context, *DetailRequest) (*DetailResponse, error)
	mustEmbedUnimplementedNachaServiceServer()
}

// UnimplementedNachaServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedNachaServiceServer struct{}

func (UnimplementedNachaServiceServer) ValidateFile(context.Context, *FileRequest) (*ValidationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateFile not implemented")
}
func (UnimplementedNachaServiceServer) CreateFile(context.Context, *NachaFileRequest) (*FileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateFile not implemented")
}
func (UnimplementedNachaServiceServer) ExportFile(context.Context, *ExportRequest) (*ExportResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExportFile not implemented")
}
func (UnimplementedNachaServiceServer) ViewFile(context.Context, *FileRequest) (*FileDetailsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ViewFile not implemented")
}
func (UnimplementedNachaServiceServer) ViewDetails(context.Context, *DetailRequest) (*DetailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ViewDetails not implemented")
}
func (UnimplementedNachaServiceServer) mustEmbedUnimplementedNachaServiceServer() {}
func (UnimplementedNachaServiceServer) testEmbeddedByValue()                      {}

// UnsafeNachaServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NachaServiceServer will
// result in compilation errors.
type UnsafeNachaServiceServer interface {
	mustEmbedUnimplementedNachaServiceServer()
}

func RegisterNachaServiceServer(s grpc.ServiceRegistrar, srv NachaServiceServer) {
	// If the following call pancis, it indicates UnimplementedNachaServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&NachaService_ServiceDesc, srv)
}

func _NachaService_ValidateFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NachaServiceServer).ValidateFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NachaService_ValidateFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NachaServiceServer).ValidateFile(ctx, req.(*FileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NachaService_CreateFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NachaFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NachaServiceServer).CreateFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NachaService_CreateFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NachaServiceServer).CreateFile(ctx, req.(*NachaFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NachaService_ExportFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NachaServiceServer).ExportFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NachaService_ExportFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NachaServiceServer).ExportFile(ctx, req.(*ExportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NachaService_ViewFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NachaServiceServer).ViewFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NachaService_ViewFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NachaServiceServer).ViewFile(ctx, req.(*FileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NachaService_ViewDetails_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DetailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NachaServiceServer).ViewDetails(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NachaService_ViewDetails_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NachaServiceServer).ViewDetails(ctx, req.(*DetailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// NachaService_ServiceDesc is the grpc.ServiceDesc for NachaService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NachaService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "nacha.NachaService",
	HandlerType: (*NachaServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ValidateFile",
			Handler:    _NachaService_ValidateFile_Handler,
		},
		{
			MethodName: "CreateFile",
			Handler:    _NachaService_CreateFile_Handler,
		},
		{
			MethodName: "ExportFile",
			Handler:    _NachaService_ExportFile_Handler,
		},
		{
			MethodName: "ViewFile",
			Handler:    _NachaService_ViewFile_Handler,
		},
		{
			MethodName: "ViewDetails",
			Handler:    _NachaService_ViewDetails_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/nacha.proto",
}
