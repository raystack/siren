// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package siren

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

// SirenServiceClient is the client API for SirenService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SirenServiceClient interface {
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	GetAlertHistory(ctx context.Context, in *GetAlertHistoryRequest, opts ...grpc.CallOption) (*GetAlertHistoryResponse, error)
	CreateAlertHistory(ctx context.Context, in *CreateAlertHistoryRequest, opts ...grpc.CallOption) (*CreateAlertHistoryResponse, error)
	GetWorkspaceChannels(ctx context.Context, in *GetWorkspaceChannelsRequest, opts ...grpc.CallOption) (*GetWorkspaceChannelsResponse, error)
	ExchangeCode(ctx context.Context, in *ExchangeCodeRequest, opts ...grpc.CallOption) (*ExchangeCodeResponse, error)
}

type sirenServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSirenServiceClient(cc grpc.ClientConnInterface) SirenServiceClient {
	return &sirenServiceClient{cc}
}

func (c *sirenServiceClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/odpf.siren.SirenService/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sirenServiceClient) GetAlertHistory(ctx context.Context, in *GetAlertHistoryRequest, opts ...grpc.CallOption) (*GetAlertHistoryResponse, error) {
	out := new(GetAlertHistoryResponse)
	err := c.cc.Invoke(ctx, "/odpf.siren.SirenService/GetAlertHistory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sirenServiceClient) CreateAlertHistory(ctx context.Context, in *CreateAlertHistoryRequest, opts ...grpc.CallOption) (*CreateAlertHistoryResponse, error) {
	out := new(CreateAlertHistoryResponse)
	err := c.cc.Invoke(ctx, "/odpf.siren.SirenService/CreateAlertHistory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sirenServiceClient) GetWorkspaceChannels(ctx context.Context, in *GetWorkspaceChannelsRequest, opts ...grpc.CallOption) (*GetWorkspaceChannelsResponse, error) {
	out := new(GetWorkspaceChannelsResponse)
	err := c.cc.Invoke(ctx, "/odpf.siren.SirenService/GetWorkspaceChannels", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sirenServiceClient) ExchangeCode(ctx context.Context, in *ExchangeCodeRequest, opts ...grpc.CallOption) (*ExchangeCodeResponse, error) {
	out := new(ExchangeCodeResponse)
	err := c.cc.Invoke(ctx, "/odpf.siren.SirenService/ExchangeCode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SirenServiceServer is the server API for SirenService service.
// All implementations must embed UnimplementedSirenServiceServer
// for forward compatibility
type SirenServiceServer interface {
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	GetAlertHistory(context.Context, *GetAlertHistoryRequest) (*GetAlertHistoryResponse, error)
	CreateAlertHistory(context.Context, *CreateAlertHistoryRequest) (*CreateAlertHistoryResponse, error)
	GetWorkspaceChannels(context.Context, *GetWorkspaceChannelsRequest) (*GetWorkspaceChannelsResponse, error)
	ExchangeCode(context.Context, *ExchangeCodeRequest) (*ExchangeCodeResponse, error)
	mustEmbedUnimplementedSirenServiceServer()
}

// UnimplementedSirenServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSirenServiceServer struct {
}

func (UnimplementedSirenServiceServer) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedSirenServiceServer) GetAlertHistory(context.Context, *GetAlertHistoryRequest) (*GetAlertHistoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAlertHistory not implemented")
}
func (UnimplementedSirenServiceServer) CreateAlertHistory(context.Context, *CreateAlertHistoryRequest) (*CreateAlertHistoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAlertHistory not implemented")
}
func (UnimplementedSirenServiceServer) GetWorkspaceChannels(context.Context, *GetWorkspaceChannelsRequest) (*GetWorkspaceChannelsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetWorkspaceChannels not implemented")
}
func (UnimplementedSirenServiceServer) ExchangeCode(context.Context, *ExchangeCodeRequest) (*ExchangeCodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExchangeCode not implemented")
}
func (UnimplementedSirenServiceServer) mustEmbedUnimplementedSirenServiceServer() {}

// UnsafeSirenServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SirenServiceServer will
// result in compilation errors.
type UnsafeSirenServiceServer interface {
	mustEmbedUnimplementedSirenServiceServer()
}

func RegisterSirenServiceServer(s grpc.ServiceRegistrar, srv SirenServiceServer) {
	s.RegisterService(&SirenService_ServiceDesc, srv)
}

func _SirenService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SirenServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/odpf.siren.SirenService/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SirenServiceServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SirenService_GetAlertHistory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAlertHistoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SirenServiceServer).GetAlertHistory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/odpf.siren.SirenService/GetAlertHistory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SirenServiceServer).GetAlertHistory(ctx, req.(*GetAlertHistoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SirenService_CreateAlertHistory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAlertHistoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SirenServiceServer).CreateAlertHistory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/odpf.siren.SirenService/CreateAlertHistory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SirenServiceServer).CreateAlertHistory(ctx, req.(*CreateAlertHistoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SirenService_GetWorkspaceChannels_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetWorkspaceChannelsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SirenServiceServer).GetWorkspaceChannels(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/odpf.siren.SirenService/GetWorkspaceChannels",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SirenServiceServer).GetWorkspaceChannels(ctx, req.(*GetWorkspaceChannelsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SirenService_ExchangeCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExchangeCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SirenServiceServer).ExchangeCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/odpf.siren.SirenService/ExchangeCode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SirenServiceServer).ExchangeCode(ctx, req.(*ExchangeCodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SirenService_ServiceDesc is the grpc.ServiceDesc for SirenService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SirenService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "odpf.siren.SirenService",
	HandlerType: (*SirenServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _SirenService_Ping_Handler,
		},
		{
			MethodName: "GetAlertHistory",
			Handler:    _SirenService_GetAlertHistory_Handler,
		},
		{
			MethodName: "CreateAlertHistory",
			Handler:    _SirenService_CreateAlertHistory_Handler,
		},
		{
			MethodName: "GetWorkspaceChannels",
			Handler:    _SirenService_GetWorkspaceChannels_Handler,
		},
		{
			MethodName: "ExchangeCode",
			Handler:    _SirenService_ExchangeCode_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "odpf/siren/siren.proto",
}
