// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: chatservice.proto

package chatpb

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

// ChatServiceClient is the client API for ChatService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChatServiceClient interface {
	Chat(ctx context.Context, opts ...grpc.CallOption) (ChatService_ChatClient, error)
	NewConversation(ctx context.Context, in *NewConversationRequest, opts ...grpc.CallOption) (*NewConversationResponse, error)
	GetConversation(ctx context.Context, in *ConversationPostRequest, opts ...grpc.CallOption) (*ConversationPostResponse, error)
	TriggerSendSMS(ctx context.Context, in *TriggerSendSMSRequest, opts ...grpc.CallOption) (*TriggerSendSMSResponse, error)
	CloseConversation(ctx context.Context, in *CloseConversationRequest, opts ...grpc.CallOption) (*CloseConversationResponse, error)
}

type chatServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewChatServiceClient(cc grpc.ClientConnInterface) ChatServiceClient {
	return &chatServiceClient{cc}
}

func (c *chatServiceClient) Chat(ctx context.Context, opts ...grpc.CallOption) (ChatService_ChatClient, error) {
	stream, err := c.cc.NewStream(ctx, &ChatService_ServiceDesc.Streams[0], "/chatservice.ChatService/Chat", opts...)
	if err != nil {
		return nil, err
	}
	x := &chatServiceChatClient{stream}
	return x, nil
}

type ChatService_ChatClient interface {
	Send(*ChatMessage) error
	Recv() (*ChatMessage, error)
	grpc.ClientStream
}

type chatServiceChatClient struct {
	grpc.ClientStream
}

func (x *chatServiceChatClient) Send(m *ChatMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *chatServiceChatClient) Recv() (*ChatMessage, error) {
	m := new(ChatMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *chatServiceClient) NewConversation(ctx context.Context, in *NewConversationRequest, opts ...grpc.CallOption) (*NewConversationResponse, error) {
	out := new(NewConversationResponse)
	err := c.cc.Invoke(ctx, "/chatservice.ChatService/NewConversation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) GetConversation(ctx context.Context, in *ConversationPostRequest, opts ...grpc.CallOption) (*ConversationPostResponse, error) {
	out := new(ConversationPostResponse)
	err := c.cc.Invoke(ctx, "/chatservice.ChatService/GetConversation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) TriggerSendSMS(ctx context.Context, in *TriggerSendSMSRequest, opts ...grpc.CallOption) (*TriggerSendSMSResponse, error) {
	out := new(TriggerSendSMSResponse)
	err := c.cc.Invoke(ctx, "/chatservice.ChatService/TriggerSendSMS", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) CloseConversation(ctx context.Context, in *CloseConversationRequest, opts ...grpc.CallOption) (*CloseConversationResponse, error) {
	out := new(CloseConversationResponse)
	err := c.cc.Invoke(ctx, "/chatservice.ChatService/CloseConversation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChatServiceServer is the server API for ChatService service.
// All implementations must embed UnimplementedChatServiceServer
// for forward compatibility
type ChatServiceServer interface {
	Chat(ChatService_ChatServer) error
	NewConversation(context.Context, *NewConversationRequest) (*NewConversationResponse, error)
	GetConversation(context.Context, *ConversationPostRequest) (*ConversationPostResponse, error)
	TriggerSendSMS(context.Context, *TriggerSendSMSRequest) (*TriggerSendSMSResponse, error)
	CloseConversation(context.Context, *CloseConversationRequest) (*CloseConversationResponse, error)
	mustEmbedUnimplementedChatServiceServer()
}

// UnimplementedChatServiceServer must be embedded to have forward compatible implementations.
type UnimplementedChatServiceServer struct {
}

func (UnimplementedChatServiceServer) Chat(ChatService_ChatServer) error {
	return status.Errorf(codes.Unimplemented, "method Chat not implemented")
}
func (UnimplementedChatServiceServer) NewConversation(context.Context, *NewConversationRequest) (*NewConversationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewConversation not implemented")
}
func (UnimplementedChatServiceServer) GetConversation(context.Context, *ConversationPostRequest) (*ConversationPostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConversation not implemented")
}
func (UnimplementedChatServiceServer) TriggerSendSMS(context.Context, *TriggerSendSMSRequest) (*TriggerSendSMSResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TriggerSendSMS not implemented")
}
func (UnimplementedChatServiceServer) CloseConversation(context.Context, *CloseConversationRequest) (*CloseConversationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CloseConversation not implemented")
}
func (UnimplementedChatServiceServer) mustEmbedUnimplementedChatServiceServer() {}

// UnsafeChatServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChatServiceServer will
// result in compilation errors.
type UnsafeChatServiceServer interface {
	mustEmbedUnimplementedChatServiceServer()
}

func RegisterChatServiceServer(s grpc.ServiceRegistrar, srv ChatServiceServer) {
	s.RegisterService(&ChatService_ServiceDesc, srv)
}

func _ChatService_Chat_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ChatServiceServer).Chat(&chatServiceChatServer{stream})
}

type ChatService_ChatServer interface {
	Send(*ChatMessage) error
	Recv() (*ChatMessage, error)
	grpc.ServerStream
}

type chatServiceChatServer struct {
	grpc.ServerStream
}

func (x *chatServiceChatServer) Send(m *ChatMessage) error {
	return x.ServerStream.SendMsg(m)
}

func (x *chatServiceChatServer) Recv() (*ChatMessage, error) {
	m := new(ChatMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _ChatService_NewConversation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewConversationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).NewConversation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chatservice.ChatService/NewConversation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).NewConversation(ctx, req.(*NewConversationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_GetConversation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConversationPostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).GetConversation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chatservice.ChatService/GetConversation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).GetConversation(ctx, req.(*ConversationPostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_TriggerSendSMS_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TriggerSendSMSRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).TriggerSendSMS(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chatservice.ChatService/TriggerSendSMS",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).TriggerSendSMS(ctx, req.(*TriggerSendSMSRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_CloseConversation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CloseConversationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).CloseConversation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chatservice.ChatService/CloseConversation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).CloseConversation(ctx, req.(*CloseConversationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ChatService_ServiceDesc is the grpc.ServiceDesc for ChatService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChatService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chatservice.ChatService",
	HandlerType: (*ChatServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NewConversation",
			Handler:    _ChatService_NewConversation_Handler,
		},
		{
			MethodName: "GetConversation",
			Handler:    _ChatService_GetConversation_Handler,
		},
		{
			MethodName: "TriggerSendSMS",
			Handler:    _ChatService_TriggerSendSMS_Handler,
		},
		{
			MethodName: "CloseConversation",
			Handler:    _ChatService_CloseConversation_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Chat",
			Handler:       _ChatService_Chat_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "chatservice.proto",
}

// ChatHTTPClient is the client API for ChatHTTP service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChatHTTPClient interface {
	FetchPost(ctx context.Context, in *FetchPostRequest, opts ...grpc.CallOption) (*FetchPostResponse, error)
	FetchOnePost(ctx context.Context, in *FetchOnePostRequest, opts ...grpc.CallOption) (*FetchOnePostResponse, error)
	SendPost(ctx context.Context, in *SendPostRequest, opts ...grpc.CallOption) (*SendPostResponse, error)
	ConversationPost(ctx context.Context, in *ConversationPostRequest, opts ...grpc.CallOption) (*ConversationPostResponse, error)
	NotificationGet(ctx context.Context, in *NotificationGetRequest, opts ...grpc.CallOption) (*NotificationGetResponse, error)
	NotificationSeenPut(ctx context.Context, in *NotificationSeenPutRequest, opts ...grpc.CallOption) (*NotificationSeenPutResponse, error)
}

type chatHTTPClient struct {
	cc grpc.ClientConnInterface
}

func NewChatHTTPClient(cc grpc.ClientConnInterface) ChatHTTPClient {
	return &chatHTTPClient{cc}
}

func (c *chatHTTPClient) FetchPost(ctx context.Context, in *FetchPostRequest, opts ...grpc.CallOption) (*FetchPostResponse, error) {
	out := new(FetchPostResponse)
	err := c.cc.Invoke(ctx, "/chatservice.ChatHTTP/FetchPost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatHTTPClient) FetchOnePost(ctx context.Context, in *FetchOnePostRequest, opts ...grpc.CallOption) (*FetchOnePostResponse, error) {
	out := new(FetchOnePostResponse)
	err := c.cc.Invoke(ctx, "/chatservice.ChatHTTP/FetchOnePost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatHTTPClient) SendPost(ctx context.Context, in *SendPostRequest, opts ...grpc.CallOption) (*SendPostResponse, error) {
	out := new(SendPostResponse)
	err := c.cc.Invoke(ctx, "/chatservice.ChatHTTP/SendPost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatHTTPClient) ConversationPost(ctx context.Context, in *ConversationPostRequest, opts ...grpc.CallOption) (*ConversationPostResponse, error) {
	out := new(ConversationPostResponse)
	err := c.cc.Invoke(ctx, "/chatservice.ChatHTTP/ConversationPost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatHTTPClient) NotificationGet(ctx context.Context, in *NotificationGetRequest, opts ...grpc.CallOption) (*NotificationGetResponse, error) {
	out := new(NotificationGetResponse)
	err := c.cc.Invoke(ctx, "/chatservice.ChatHTTP/NotificationGet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatHTTPClient) NotificationSeenPut(ctx context.Context, in *NotificationSeenPutRequest, opts ...grpc.CallOption) (*NotificationSeenPutResponse, error) {
	out := new(NotificationSeenPutResponse)
	err := c.cc.Invoke(ctx, "/chatservice.ChatHTTP/NotificationSeenPut", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChatHTTPServer is the server API for ChatHTTP service.
// All implementations must embed UnimplementedChatHTTPServer
// for forward compatibility
type ChatHTTPServer interface {
	FetchPost(context.Context, *FetchPostRequest) (*FetchPostResponse, error)
	FetchOnePost(context.Context, *FetchOnePostRequest) (*FetchOnePostResponse, error)
	SendPost(context.Context, *SendPostRequest) (*SendPostResponse, error)
	ConversationPost(context.Context, *ConversationPostRequest) (*ConversationPostResponse, error)
	NotificationGet(context.Context, *NotificationGetRequest) (*NotificationGetResponse, error)
	NotificationSeenPut(context.Context, *NotificationSeenPutRequest) (*NotificationSeenPutResponse, error)
	mustEmbedUnimplementedChatHTTPServer()
}

// UnimplementedChatHTTPServer must be embedded to have forward compatible implementations.
type UnimplementedChatHTTPServer struct {
}

func (UnimplementedChatHTTPServer) FetchPost(context.Context, *FetchPostRequest) (*FetchPostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FetchPost not implemented")
}
func (UnimplementedChatHTTPServer) FetchOnePost(context.Context, *FetchOnePostRequest) (*FetchOnePostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FetchOnePost not implemented")
}
func (UnimplementedChatHTTPServer) SendPost(context.Context, *SendPostRequest) (*SendPostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendPost not implemented")
}
func (UnimplementedChatHTTPServer) ConversationPost(context.Context, *ConversationPostRequest) (*ConversationPostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConversationPost not implemented")
}
func (UnimplementedChatHTTPServer) NotificationGet(context.Context, *NotificationGetRequest) (*NotificationGetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NotificationGet not implemented")
}
func (UnimplementedChatHTTPServer) NotificationSeenPut(context.Context, *NotificationSeenPutRequest) (*NotificationSeenPutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NotificationSeenPut not implemented")
}
func (UnimplementedChatHTTPServer) mustEmbedUnimplementedChatHTTPServer() {}

// UnsafeChatHTTPServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChatHTTPServer will
// result in compilation errors.
type UnsafeChatHTTPServer interface {
	mustEmbedUnimplementedChatHTTPServer()
}

func RegisterChatHTTPServer(s grpc.ServiceRegistrar, srv ChatHTTPServer) {
	s.RegisterService(&ChatHTTP_ServiceDesc, srv)
}

func _ChatHTTP_FetchPost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FetchPostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatHTTPServer).FetchPost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chatservice.ChatHTTP/FetchPost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatHTTPServer).FetchPost(ctx, req.(*FetchPostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatHTTP_FetchOnePost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FetchOnePostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatHTTPServer).FetchOnePost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chatservice.ChatHTTP/FetchOnePost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatHTTPServer).FetchOnePost(ctx, req.(*FetchOnePostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatHTTP_SendPost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendPostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatHTTPServer).SendPost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chatservice.ChatHTTP/SendPost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatHTTPServer).SendPost(ctx, req.(*SendPostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatHTTP_ConversationPost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConversationPostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatHTTPServer).ConversationPost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chatservice.ChatHTTP/ConversationPost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatHTTPServer).ConversationPost(ctx, req.(*ConversationPostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatHTTP_NotificationGet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NotificationGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatHTTPServer).NotificationGet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chatservice.ChatHTTP/NotificationGet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatHTTPServer).NotificationGet(ctx, req.(*NotificationGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatHTTP_NotificationSeenPut_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NotificationSeenPutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatHTTPServer).NotificationSeenPut(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chatservice.ChatHTTP/NotificationSeenPut",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatHTTPServer).NotificationSeenPut(ctx, req.(*NotificationSeenPutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ChatHTTP_ServiceDesc is the grpc.ServiceDesc for ChatHTTP service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChatHTTP_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chatservice.ChatHTTP",
	HandlerType: (*ChatHTTPServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FetchPost",
			Handler:    _ChatHTTP_FetchPost_Handler,
		},
		{
			MethodName: "FetchOnePost",
			Handler:    _ChatHTTP_FetchOnePost_Handler,
		},
		{
			MethodName: "SendPost",
			Handler:    _ChatHTTP_SendPost_Handler,
		},
		{
			MethodName: "ConversationPost",
			Handler:    _ChatHTTP_ConversationPost_Handler,
		},
		{
			MethodName: "NotificationGet",
			Handler:    _ChatHTTP_NotificationGet_Handler,
		},
		{
			MethodName: "NotificationSeenPut",
			Handler:    _ChatHTTP_NotificationSeenPut_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "chatservice.proto",
}
