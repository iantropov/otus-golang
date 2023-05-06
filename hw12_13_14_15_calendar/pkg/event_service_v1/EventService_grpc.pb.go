// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: EventService.proto

package event_service_v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// EventServiceV1Client is the client API for EventServiceV1 service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EventServiceV1Client interface {
	// * Создать (событие);
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// * Обновить (ID события, событие);
	Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// * Удалить (ID события);
	Delete(ctx context.Context, in *IDRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// * Получить (ID события);
	Get(ctx context.Context, in *IDRequest, opts ...grpc.CallOption) (*GetResponse, error)
	// * СписокСобытийНаДень (дата);
	ListEventForDay(ctx context.Context, in *TimeRequest, opts ...grpc.CallOption) (*ListEventResponse, error)
	// * СписокСобытийНаНеделю (дата начала недели);
	ListEventForWeek(ctx context.Context, in *TimeRequest, opts ...grpc.CallOption) (*ListEventResponse, error)
	// * СписокСобытийНaМесяц (дата начала месяца).
	ListEventForMonth(ctx context.Context, in *TimeRequest, opts ...grpc.CallOption) (*ListEventResponse, error)
}

type eventServiceV1Client struct {
	cc grpc.ClientConnInterface
}

func NewEventServiceV1Client(cc grpc.ClientConnInterface) EventServiceV1Client {
	return &eventServiceV1Client{cc}
}

func (c *eventServiceV1Client) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/event_service_v1.EventServiceV1/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceV1Client) Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/event_service_v1.EventServiceV1/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceV1Client) Delete(ctx context.Context, in *IDRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/event_service_v1.EventServiceV1/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceV1Client) Get(ctx context.Context, in *IDRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, "/event_service_v1.EventServiceV1/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceV1Client) ListEventForDay(ctx context.Context, in *TimeRequest, opts ...grpc.CallOption) (*ListEventResponse, error) {
	out := new(ListEventResponse)
	err := c.cc.Invoke(ctx, "/event_service_v1.EventServiceV1/ListEventForDay", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceV1Client) ListEventForWeek(ctx context.Context, in *TimeRequest, opts ...grpc.CallOption) (*ListEventResponse, error) {
	out := new(ListEventResponse)
	err := c.cc.Invoke(ctx, "/event_service_v1.EventServiceV1/ListEventForWeek", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceV1Client) ListEventForMonth(ctx context.Context, in *TimeRequest, opts ...grpc.CallOption) (*ListEventResponse, error) {
	out := new(ListEventResponse)
	err := c.cc.Invoke(ctx, "/event_service_v1.EventServiceV1/ListEventForMonth", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EventServiceV1Server is the server API for EventServiceV1 service.
// All implementations must embed UnimplementedEventServiceV1Server
// for forward compatibility
type EventServiceV1Server interface {
	// * Создать (событие);
	Create(context.Context, *CreateRequest) (*emptypb.Empty, error)
	// * Обновить (ID события, событие);
	Update(context.Context, *UpdateRequest) (*emptypb.Empty, error)
	// * Удалить (ID события);
	Delete(context.Context, *IDRequest) (*emptypb.Empty, error)
	// * Получить (ID события);
	Get(context.Context, *IDRequest) (*GetResponse, error)
	// * СписокСобытийНаДень (дата);
	ListEventForDay(context.Context, *TimeRequest) (*ListEventResponse, error)
	// * СписокСобытийНаНеделю (дата начала недели);
	ListEventForWeek(context.Context, *TimeRequest) (*ListEventResponse, error)
	// * СписокСобытийНaМесяц (дата начала месяца).
	ListEventForMonth(context.Context, *TimeRequest) (*ListEventResponse, error)
	mustEmbedUnimplementedEventServiceV1Server()
}

// UnimplementedEventServiceV1Server must be embedded to have forward compatible implementations.
type UnimplementedEventServiceV1Server struct {
}

func (UnimplementedEventServiceV1Server) Create(context.Context, *CreateRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedEventServiceV1Server) Update(context.Context, *UpdateRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedEventServiceV1Server) Delete(context.Context, *IDRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedEventServiceV1Server) Get(context.Context, *IDRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedEventServiceV1Server) ListEventForDay(context.Context, *TimeRequest) (*ListEventResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEventForDay not implemented")
}
func (UnimplementedEventServiceV1Server) ListEventForWeek(context.Context, *TimeRequest) (*ListEventResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEventForWeek not implemented")
}
func (UnimplementedEventServiceV1Server) ListEventForMonth(context.Context, *TimeRequest) (*ListEventResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEventForMonth not implemented")
}
func (UnimplementedEventServiceV1Server) mustEmbedUnimplementedEventServiceV1Server() {}

// UnsafeEventServiceV1Server may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EventServiceV1Server will
// result in compilation errors.
type UnsafeEventServiceV1Server interface {
	mustEmbedUnimplementedEventServiceV1Server()
}

func RegisterEventServiceV1Server(s grpc.ServiceRegistrar, srv EventServiceV1Server) {
	s.RegisterService(&EventServiceV1_ServiceDesc, srv)
}

func _EventServiceV1_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceV1Server).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event_service_v1.EventServiceV1/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceV1Server).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventServiceV1_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceV1Server).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event_service_v1.EventServiceV1/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceV1Server).Update(ctx, req.(*UpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventServiceV1_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceV1Server).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event_service_v1.EventServiceV1/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceV1Server).Delete(ctx, req.(*IDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventServiceV1_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceV1Server).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event_service_v1.EventServiceV1/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceV1Server).Get(ctx, req.(*IDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventServiceV1_ListEventForDay_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TimeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceV1Server).ListEventForDay(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event_service_v1.EventServiceV1/ListEventForDay",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceV1Server).ListEventForDay(ctx, req.(*TimeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventServiceV1_ListEventForWeek_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TimeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceV1Server).ListEventForWeek(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event_service_v1.EventServiceV1/ListEventForWeek",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceV1Server).ListEventForWeek(ctx, req.(*TimeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventServiceV1_ListEventForMonth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TimeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceV1Server).ListEventForMonth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event_service_v1.EventServiceV1/ListEventForMonth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceV1Server).ListEventForMonth(ctx, req.(*TimeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// EventServiceV1_ServiceDesc is the grpc.ServiceDesc for EventServiceV1 service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EventServiceV1_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "event_service_v1.EventServiceV1",
	HandlerType: (*EventServiceV1Server)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _EventServiceV1_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _EventServiceV1_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _EventServiceV1_Delete_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _EventServiceV1_Get_Handler,
		},
		{
			MethodName: "ListEventForDay",
			Handler:    _EventServiceV1_ListEventForDay_Handler,
		},
		{
			MethodName: "ListEventForWeek",
			Handler:    _EventServiceV1_ListEventForWeek_Handler,
		},
		{
			MethodName: "ListEventForMonth",
			Handler:    _EventServiceV1_ListEventForMonth_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "EventService.proto",
}
