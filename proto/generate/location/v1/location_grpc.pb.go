// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.6.1
// source: proto/location/location.proto

package v1

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

const (
	Courier_GetCourierLatestPosition_FullMethodName = "/Courier/GetCourierLatestPosition"
)

// CourierClient is the client API for Courier service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CourierClient interface {
	GetCourierLatestPosition(ctx context.Context, in *GetCourierLatestPositionRequest, opts ...grpc.CallOption) (*GetCourierLatestPositionResponse, error)
}

type courierClient struct {
	cc grpc.ClientConnInterface
}

func NewCourierClient(cc grpc.ClientConnInterface) CourierClient {
	return &courierClient{cc}
}

func (c *courierClient) GetCourierLatestPosition(ctx context.Context, in *GetCourierLatestPositionRequest, opts ...grpc.CallOption) (*GetCourierLatestPositionResponse, error) {
	out := new(GetCourierLatestPositionResponse)
	err := c.cc.Invoke(ctx, Courier_GetCourierLatestPosition_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CourierServer is the server API for Courier service.
// All implementations must embed UnimplementedCourierServer
// for forward compatibility
type CourierServer interface {
	GetCourierLatestPosition(context.Context, *GetCourierLatestPositionRequest) (*GetCourierLatestPositionResponse, error)
	mustEmbedUnimplementedCourierServer()
}

// UnimplementedCourierServer must be embedded to have forward compatible implementations.
type UnimplementedCourierServer struct {
}

func (UnimplementedCourierServer) GetCourierLatestPosition(context.Context, *GetCourierLatestPositionRequest) (*GetCourierLatestPositionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCourierLatestPosition not implemented")
}
func (UnimplementedCourierServer) mustEmbedUnimplementedCourierServer() {}

// UnsafeCourierServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CourierServer will
// result in compilation errors.
type UnsafeCourierServer interface {
	mustEmbedUnimplementedCourierServer()
}

func RegisterCourierServer(s grpc.ServiceRegistrar, srv CourierServer) {
	s.RegisterService(&Courier_ServiceDesc, srv)
}

func _Courier_GetCourierLatestPosition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCourierLatestPositionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CourierServer).GetCourierLatestPosition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Courier_GetCourierLatestPosition_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CourierServer).GetCourierLatestPosition(ctx, req.(*GetCourierLatestPositionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Courier_ServiceDesc is the grpc.ServiceDesc for Courier service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Courier_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Courier",
	HandlerType: (*CourierServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCourierLatestPosition",
			Handler:    _Courier_GetCourierLatestPosition_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/location/location.proto",
}