// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: post_storage.proto

package post_storage

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
	PostStorage_GetPost_FullMethodName = "/PostStorage/GetPost"
)

// PostStorageClient is the client API for PostStorage service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PostStorageClient interface {
	GetPost(ctx context.Context, in *GetPostRequest, opts ...grpc.CallOption) (*GetPostResponse, error)
}

type postStorageClient struct {
	cc grpc.ClientConnInterface
}

func NewPostStorageClient(cc grpc.ClientConnInterface) PostStorageClient {
	return &postStorageClient{cc}
}

func (c *postStorageClient) GetPost(ctx context.Context, in *GetPostRequest, opts ...grpc.CallOption) (*GetPostResponse, error) {
	out := new(GetPostResponse)
	err := c.cc.Invoke(ctx, PostStorage_GetPost_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PostStorageServer is the server API for PostStorage service.
// All implementations must embed UnimplementedPostStorageServer
// for forward compatibility
type PostStorageServer interface {
	GetPost(context.Context, *GetPostRequest) (*GetPostResponse, error)
	mustEmbedUnimplementedPostStorageServer()
}

// UnimplementedPostStorageServer must be embedded to have forward compatible implementations.
type UnimplementedPostStorageServer struct {
}

func (UnimplementedPostStorageServer) GetPost(context.Context, *GetPostRequest) (*GetPostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPost not implemented")
}
func (UnimplementedPostStorageServer) mustEmbedUnimplementedPostStorageServer() {}

// UnsafePostStorageServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PostStorageServer will
// result in compilation errors.
type UnsafePostStorageServer interface {
	mustEmbedUnimplementedPostStorageServer()
}

func RegisterPostStorageServer(s grpc.ServiceRegistrar, srv PostStorageServer) {
	s.RegisterService(&PostStorage_ServiceDesc, srv)
}

func _PostStorage_GetPost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostStorageServer).GetPost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostStorage_GetPost_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostStorageServer).GetPost(ctx, req.(*GetPostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PostStorage_ServiceDesc is the grpc.ServiceDesc for PostStorage service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PostStorage_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "PostStorage",
	HandlerType: (*PostStorageServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPost",
			Handler:    _PostStorage_GetPost_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "post_storage.proto",
}
