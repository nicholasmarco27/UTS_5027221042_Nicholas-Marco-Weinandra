// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.26.1
// source: menulist.proto

package menulist

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MenuApiClient is the client API for MenuApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MenuApiClient interface {
	CreateMenu(ctx context.Context, in *Menu, opts ...grpc.CallOption) (*Menu, error)
	ListMenus(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*MenuList, error)
	UpdateMenu(ctx context.Context, in *Menu, opts ...grpc.CallOption) (*Menu, error)
	DeleteMenu(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (*wrapperspb.BoolValue, error)
}

type menuApiClient struct {
	cc grpc.ClientConnInterface
}

func NewMenuApiClient(cc grpc.ClientConnInterface) MenuApiClient {
	return &menuApiClient{cc}
}

func (c *menuApiClient) CreateMenu(ctx context.Context, in *Menu, opts ...grpc.CallOption) (*Menu, error) {
	out := new(Menu)
	err := c.cc.Invoke(ctx, "/protoapi.MenuApi/CreateMenu", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *menuApiClient) ListMenus(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*MenuList, error) {
	out := new(MenuList)
	err := c.cc.Invoke(ctx, "/protoapi.MenuApi/ListMenus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *menuApiClient) UpdateMenu(ctx context.Context, in *Menu, opts ...grpc.CallOption) (*Menu, error) {
	out := new(Menu)
	err := c.cc.Invoke(ctx, "/protoapi.MenuApi/UpdateMenu", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *menuApiClient) DeleteMenu(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (*wrapperspb.BoolValue, error) {
	out := new(wrapperspb.BoolValue)
	err := c.cc.Invoke(ctx, "/protoapi.MenuApi/DeleteMenu", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MenuApiServer is the server API for MenuApi service.
// All implementations must embed UnimplementedMenuApiServer
// for forward compatibility
type MenuApiServer interface {
	CreateMenu(context.Context, *Menu) (*Menu, error)
	ListMenus(context.Context, *emptypb.Empty) (*MenuList, error)
	UpdateMenu(context.Context, *Menu) (*Menu, error)
	DeleteMenu(context.Context, *wrapperspb.StringValue) (*wrapperspb.BoolValue, error)
}

// UnimplementedMenuApiServer must be embedded to have forward compatible implementations.
type UnimplementedMenuApiServer struct {
}

func (UnimplementedMenuApiServer) CreateMenu(context.Context, *Menu) (*Menu, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMenu not implemented")
}
func (UnimplementedMenuApiServer) ListMenus(context.Context, *emptypb.Empty) (*MenuList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListMenus not implemented")
}
func (UnimplementedMenuApiServer) UpdateMenu(context.Context, *Menu) (*Menu, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMenu not implemented")
}
func (UnimplementedMenuApiServer) DeleteMenu(context.Context, *wrapperspb.StringValue) (*wrapperspb.BoolValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMenu not implemented")
}

func RegisterMenuApiServer(s grpc.ServiceRegistrar, srv MenuApiServer) {
	s.RegisterService(&MenuApi_ServiceDesc, srv)
}

func _MenuApi_CreateMenu_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Menu)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MenuApiServer).CreateMenu(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protoapi.MenuApi/CreateMenu",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MenuApiServer).CreateMenu(ctx, req.(*Menu))
	}
	return interceptor(ctx, in, info, handler)
}

func _MenuApi_ListMenus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MenuApiServer).ListMenus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protoapi.MenuApi/ListMenus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MenuApiServer).ListMenus(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _MenuApi_UpdateMenu_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Menu)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MenuApiServer).UpdateMenu(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protoapi.MenuApi/UpdateMenu",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MenuApiServer).UpdateMenu(ctx, req.(*Menu))
	}
	return interceptor(ctx, in, info, handler)
}

func _MenuApi_DeleteMenu_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(wrapperspb.StringValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MenuApiServer).DeleteMenu(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protoapi.MenuApi/DeleteMenu",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MenuApiServer).DeleteMenu(ctx, req.(*wrapperspb.StringValue))
	}
	return interceptor(ctx, in, info, handler)
}

// MenuApi_ServiceDesc is the grpc.ServiceDesc for MenuApi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MenuApi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protoapi.MenuApi",
	HandlerType: (*MenuApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateMenu",
			Handler:    _MenuApi_CreateMenu_Handler,
		},
		{
			MethodName: "ListMenus",
			Handler:    _MenuApi_ListMenus_Handler,
		},
		{
			MethodName: "UpdateMenu",
			Handler:    _MenuApi_UpdateMenu_Handler,
		},
		{
			MethodName: "DeleteMenu",
			Handler:    _MenuApi_DeleteMenu_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "menulist.proto",
}
