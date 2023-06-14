// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package testpb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// IOServiceClient is the client API for IOService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type IOServiceClient interface {
	Call(ctx context.Context, in *Input, opts ...grpc.CallOption) (*Output, error)
}

type iOServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewIOServiceClient(cc grpc.ClientConnInterface) IOServiceClient {
	return &iOServiceClient{cc}
}

func (c *iOServiceClient) Call(ctx context.Context, in *Input, opts ...grpc.CallOption) (*Output, error) {
	out := new(Output)
	err := c.cc.Invoke(ctx, "/IOService/Call", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IOServiceServer is the server API for IOService service.
// All implementations must embed UnimplementedIOServiceServer
// for forward compatibility
type IOServiceServer interface {
	Call(context.Context, *Input) (*Output, error)
	mustEmbedUnimplementedIOServiceServer()
}

// UnimplementedIOServiceServer must be embedded to have forward compatible implementations.
type UnimplementedIOServiceServer struct {
}

func (UnimplementedIOServiceServer) Call(context.Context, *Input) (*Output, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Call not implemented")
}
func (UnimplementedIOServiceServer) mustEmbedUnimplementedIOServiceServer() {}

// UnsafeIOServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to IOServiceServer will
// result in compilation errors.
type UnsafeIOServiceServer interface {
	mustEmbedUnimplementedIOServiceServer()
}

func RegisterIOServiceServer(s grpc.ServiceRegistrar, srv IOServiceServer) {
	s.RegisterService(&IOService_ServiceDesc, srv)
}

func _IOService_Call_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Input)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IOServiceServer).Call(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/IOService/Call",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IOServiceServer).Call(ctx, req.(*Input))
	}
	return interceptor(ctx, in, info, handler)
}

// IOService_ServiceDesc is the grpc.ServiceDesc for IOService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var IOService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "IOService",
	HandlerType: (*IOServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Call",
			Handler:    _IOService_Call_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "io.proto",
}
