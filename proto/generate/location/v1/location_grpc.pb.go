// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.6.1
// source: location/location.proto

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
	CourierLocationPosition_GetCourierLatestPosition_FullMethodName = "/CourierLocationPosition/GetCourierLatestPosition"
)

// CourierLocationPositionClient is the client API for CourierLocationPosition service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CourierLocationPositionClient interface {
	GetCourierLatestPosition(ctx context.Context, in *GetCourierLatestPositionRequest, opts ...grpc.CallOption) (*GetCourierLatestPositionResponse, error)
}

type courierLocationPositionClient struct {
	cc grpc.ClientConnInterface
}

func NewCourierLocationPositionClient(cc grpc.ClientConnInterface) CourierLocationPositionClient {
	return &courierLocationPositionClient{cc}
}

func (c *courierLocationPositionClient) GetCourierLatestPosition(ctx context.Context, in *GetCourierLatestPositionRequest, opts ...grpc.CallOption) (*GetCourierLatestPositionResponse, error) {
	out := new(GetCourierLatestPositionResponse)
	err := c.cc.Invoke(ctx, CourierLocationPosition_GetCourierLatestPosition_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CourierLocationPositionServer is the server API for CourierLocationPosition service.
// All implementations must embed UnimplementedCourierLocationPositionServer
// for forward compatibility
type CourierLocationPositionServer interface {
	GetCourierLatestPosition(context.Context, *GetCourierLatestPositionRequest) (*GetCourierLatestPositionResponse, error)
	mustEmbedUnimplementedCourierLocationPositionServer()
}

// UnimplementedCourierLocationPositionServer must be embedded to have forward compatible implementations.
type UnimplementedCourierLocationPositionServer struct {
}

func (UnimplementedCourierLocationPositionServer) GetCourierLatestPosition(context.Context, *GetCourierLatestPositionRequest) (*GetCourierLatestPositionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCourierLatestPosition not implemented")
}
func (UnimplementedCourierLocationPositionServer) mustEmbedUnimplementedCourierLocationPositionServer() {
}

// UnsafeCourierLocationPositionServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CourierLocationPositionServer will
// result in compilation errors.
type UnsafeCourierLocationPositionServer interface {
	mustEmbedUnimplementedCourierLocationPositionServer()
}

func RegisterCourierLocationPositionServer(s grpc.ServiceRegistrar, srv CourierLocationPositionServer) {
	s.RegisterService(&CourierLocationPosition_ServiceDesc, srv)
}

func _CourierLocationPosition_GetCourierLatestPosition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCourierLatestPositionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CourierLocationPositionServer).GetCourierLatestPosition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CourierLocationPosition_GetCourierLatestPosition_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CourierLocationPositionServer).GetCourierLatestPosition(ctx, req.(*GetCourierLatestPositionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CourierLocationPosition_ServiceDesc is the grpc.ServiceDesc for CourierLocationPosition service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CourierLocationPosition_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "CourierLocationPosition",
	HandlerType: (*CourierLocationPositionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCourierLatestPosition",
			Handler:    _CourierLocationPosition_GetCourierLatestPosition_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "location/location.proto",
}