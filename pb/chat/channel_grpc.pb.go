// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v3.14.0
// source: chat/channel.proto

package chat

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
	ChannelService_SearchChannelByName_FullMethodName = "/chat.ChannelService/SearchChannelByName"
)

// ChannelServiceClient is the client API for ChannelService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChannelServiceClient interface {
	SearchChannelByName(ctx context.Context, in *SearchChannelByNameRequest, opts ...grpc.CallOption) (*SearchChannelByNameResponse, error)
}

type channelServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewChannelServiceClient(cc grpc.ClientConnInterface) ChannelServiceClient {
	return &channelServiceClient{cc}
}

func (c *channelServiceClient) SearchChannelByName(ctx context.Context, in *SearchChannelByNameRequest, opts ...grpc.CallOption) (*SearchChannelByNameResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SearchChannelByNameResponse)
	err := c.cc.Invoke(ctx, ChannelService_SearchChannelByName_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChannelServiceServer is the server API for ChannelService service.
// All implementations must embed UnimplementedChannelServiceServer
// for forward compatibility
type ChannelServiceServer interface {
	SearchChannelByName(context.Context, *SearchChannelByNameRequest) (*SearchChannelByNameResponse, error)
	mustEmbedUnimplementedChannelServiceServer()
}

// UnimplementedChannelServiceServer must be embedded to have forward compatible implementations.
type UnimplementedChannelServiceServer struct {
}

func (UnimplementedChannelServiceServer) SearchChannelByName(context.Context, *SearchChannelByNameRequest) (*SearchChannelByNameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchChannelByName not implemented")
}
func (UnimplementedChannelServiceServer) mustEmbedUnimplementedChannelServiceServer() {}

// UnsafeChannelServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChannelServiceServer will
// result in compilation errors.
type UnsafeChannelServiceServer interface {
	mustEmbedUnimplementedChannelServiceServer()
}

func RegisterChannelServiceServer(s grpc.ServiceRegistrar, srv ChannelServiceServer) {
	s.RegisterService(&ChannelService_ServiceDesc, srv)
}

func _ChannelService_SearchChannelByName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchChannelByNameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChannelServiceServer).SearchChannelByName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChannelService_SearchChannelByName_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChannelServiceServer).SearchChannelByName(ctx, req.(*SearchChannelByNameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ChannelService_ServiceDesc is the grpc.ServiceDesc for ChannelService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChannelService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chat.ChannelService",
	HandlerType: (*ChannelServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SearchChannelByName",
			Handler:    _ChannelService_SearchChannelByName_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "chat/channel.proto",
}
