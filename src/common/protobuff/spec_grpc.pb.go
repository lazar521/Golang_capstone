// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v5.27.1
// source: spec.proto

package protobuff

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	LocationHistoryService_UpdateHistory_FullMethodName = "/LocationHistoryService/UpdateHistory"
)

// LocationHistoryServiceClient is the client API for LocationHistoryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LocationHistoryServiceClient interface {
	UpdateHistory(ctx context.Context, in *LocationUpdateRequest, opts ...grpc.CallOption) (*LocationUpdateReply, error)
}

type locationHistoryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLocationHistoryServiceClient(cc grpc.ClientConnInterface) LocationHistoryServiceClient {
	return &locationHistoryServiceClient{cc}
}

func (c *locationHistoryServiceClient) UpdateHistory(ctx context.Context, in *LocationUpdateRequest, opts ...grpc.CallOption) (*LocationUpdateReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LocationUpdateReply)
	err := c.cc.Invoke(ctx, LocationHistoryService_UpdateHistory_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LocationHistoryServiceServer is the server API for LocationHistoryService service.
// All implementations must embed UnimplementedLocationHistoryServiceServer
// for forward compatibility
type LocationHistoryServiceServer interface {
	UpdateHistory(context.Context, *LocationUpdateRequest) (*LocationUpdateReply, error)
	mustEmbedUnimplementedLocationHistoryServiceServer()
}

// UnimplementedLocationHistoryServiceServer must be embedded to have forward compatible implementations.
type UnimplementedLocationHistoryServiceServer struct {
}

func (UnimplementedLocationHistoryServiceServer) UpdateHistory(context.Context, *LocationUpdateRequest) (*LocationUpdateReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateHistory not implemented")
}
func (UnimplementedLocationHistoryServiceServer) mustEmbedUnimplementedLocationHistoryServiceServer() {
}

// UnsafeLocationHistoryServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LocationHistoryServiceServer will
// result in compilation errors.
type UnsafeLocationHistoryServiceServer interface {
	mustEmbedUnimplementedLocationHistoryServiceServer()
}

func RegisterLocationHistoryServiceServer(s grpc.ServiceRegistrar, srv LocationHistoryServiceServer) {
	s.RegisterService(&LocationHistoryService_ServiceDesc, srv)
}

func _LocationHistoryService_UpdateHistory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LocationUpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocationHistoryServiceServer).UpdateHistory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LocationHistoryService_UpdateHistory_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocationHistoryServiceServer).UpdateHistory(ctx, req.(*LocationUpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LocationHistoryService_ServiceDesc is the grpc.ServiceDesc for LocationHistoryService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LocationHistoryService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "LocationHistoryService",
	HandlerType: (*LocationHistoryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateHistory",
			Handler:    _LocationHistoryService_UpdateHistory_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "spec.proto",
}
