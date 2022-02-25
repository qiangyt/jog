// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.18.1
// source: api/proto/index.proto

package index

import (
	v1 "common/v1"
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// IndexClient is the client API for Index service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type IndexClient interface {
	Get(ctx context.Context, in *v1.StringId, opts ...grpc.CallOption) (*IndexData, error)
	List(ctx context.Context, in *IndexListReq, opts ...grpc.CallOption) (*IndexListResp, error)
	Create(ctx context.Context, in *IndexCreateReq, opts ...grpc.CallOption) (*v1.StringId, error)
	Update(ctx context.Context, in *IndexUpdateReq, opts ...grpc.CallOption) (*v1.StringId, error)
	Remove(ctx context.Context, in *v1.StringId, opts ...grpc.CallOption) (*IndexData, error)
}

type indexClient struct {
	cc grpc.ClientConnInterface
}

func NewIndexClient(cc grpc.ClientConnInterface) IndexClient {
	return &indexClient{cc}
}

func (c *indexClient) Get(ctx context.Context, in *v1.StringId, opts ...grpc.CallOption) (*IndexData, error) {
	out := new(IndexData)
	err := c.cc.Invoke(ctx, "/index.v1.Index/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *indexClient) List(ctx context.Context, in *IndexListReq, opts ...grpc.CallOption) (*IndexListResp, error) {
	out := new(IndexListResp)
	err := c.cc.Invoke(ctx, "/index.v1.Index/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *indexClient) Create(ctx context.Context, in *IndexCreateReq, opts ...grpc.CallOption) (*v1.StringId, error) {
	out := new(v1.StringId)
	err := c.cc.Invoke(ctx, "/index.v1.Index/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *indexClient) Update(ctx context.Context, in *IndexUpdateReq, opts ...grpc.CallOption) (*v1.StringId, error) {
	out := new(v1.StringId)
	err := c.cc.Invoke(ctx, "/index.v1.Index/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *indexClient) Remove(ctx context.Context, in *v1.StringId, opts ...grpc.CallOption) (*IndexData, error) {
	out := new(IndexData)
	err := c.cc.Invoke(ctx, "/index.v1.Index/Remove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IndexServer is the server API for Index service.
// All implementations must embed UnimplementedIndexServer
// for forward compatibility
type IndexServer interface {
	Get(context.Context, *v1.StringId) (*IndexData, error)
	List(context.Context, *IndexListReq) (*IndexListResp, error)
	Create(context.Context, *IndexCreateReq) (*v1.StringId, error)
	Update(context.Context, *IndexUpdateReq) (*v1.StringId, error)
	Remove(context.Context, *v1.StringId) (*IndexData, error)
	mustEmbedUnimplementedIndexServer()
}

// UnimplementedIndexServer must be embedded to have forward compatible implementations.
type UnimplementedIndexServer struct {
}

func (UnimplementedIndexServer) Get(context.Context, *v1.StringId) (*IndexData, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedIndexServer) List(context.Context, *IndexListReq) (*IndexListResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedIndexServer) Create(context.Context, *IndexCreateReq) (*v1.StringId, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedIndexServer) Update(context.Context, *IndexUpdateReq) (*v1.StringId, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedIndexServer) Remove(context.Context, *v1.StringId) (*IndexData, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Remove not implemented")
}
func (UnimplementedIndexServer) mustEmbedUnimplementedIndexServer() {}

// UnsafeIndexServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to IndexServer will
// result in compilation errors.
type UnsafeIndexServer interface {
	mustEmbedUnimplementedIndexServer()
}

func RegisterIndexServer(s grpc.ServiceRegistrar, srv IndexServer) {
	s.RegisterService(&Index_ServiceDesc, srv)
}

func _Index_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.StringId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/index.v1.Index/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndexServer).Get(ctx, req.(*v1.StringId))
	}
	return interceptor(ctx, in, info, handler)
}

func _Index_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IndexListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/index.v1.Index/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndexServer).List(ctx, req.(*IndexListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Index_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IndexCreateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/index.v1.Index/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndexServer).Create(ctx, req.(*IndexCreateReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Index_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IndexUpdateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/index.v1.Index/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndexServer).Update(ctx, req.(*IndexUpdateReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Index_Remove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.StringId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexServer).Remove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/index.v1.Index/Remove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndexServer).Remove(ctx, req.(*v1.StringId))
	}
	return interceptor(ctx, in, info, handler)
}

// Index_ServiceDesc is the grpc.ServiceDesc for Index service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Index_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "index.v1.Index",
	HandlerType: (*IndexServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _Index_Get_Handler,
		},
		{
			MethodName: "List",
			Handler:    _Index_List_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _Index_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _Index_Update_Handler,
		},
		{
			MethodName: "Remove",
			Handler:    _Index_Remove_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/index.proto",
}