package chatservice

import (
	"context"
	"time"

	"github.com/aqaurius6666/cronjob/src/pb/chatpb"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ChatserviceAddr string

var (
	_       Service = (*ServiceGRPC)(nil)
	timeout         = 5 * time.Second
)

type ServiceGRPC struct {
	Ctx    context.Context
	Client chatpb.ChatServiceClient
}

func ConnectClient(ctx context.Context, addr ChatserviceAddr) (chatpb.ChatServiceClient, error) {
	nctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	conn, err := grpc.DialContext(nctx, string(addr), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithChainUnaryInterceptor(
		otelgrpc.UnaryClientInterceptor(),
	))
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}
	return chatpb.NewChatServiceClient(conn), nil
}

// Implements TriggerSendSMS
func (s *ServiceGRPC) TriggerSendSMS(ctx context.Context, userId string) error {
	nctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	_, err := s.Client.TriggerSendSMS(nctx, &chatpb.TriggerSendSMSRequest{
		UserId: userId,
	})
	if err != nil {
		return xerrors.Errorf("%w", err)
	}
	return nil
}

// func (s *ServiceGRPC) NewConversation(ctx context.Context, order string, member []string, phones []string) (string, error) {
// 	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.NewConversation))
// 	defer span.End()

// 	ctx, cancel := context.WithTimeout(ctx, timeout)
// 	defer cancel()

// 	res, err := s.Client.NewConversation(ctx, &chatpb.NewConversationRequest{
// 		OrderId:      order,
// 		MemberIds:    member,
// 		PhoneNumbers: phones,
// 	})
// 	if err != nil {
// 		err = xerrors.Errorf("%w", err)
// 		lib.RecordError(span, err, ctx)
// 		return "", err
// 	}
// 	return res.Id, nil
// }

// func (s *ServiceGRPC) GetConversation(ctx context.Context, _userId string, conversationIds []string) (*chatpb.ConversationPostResponse_Data, error) {
// 	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetConversation))
// 	defer span.End()

// 	ctx, cancel := context.WithTimeout(ctx, timeout)
// 	defer cancel()

// 	res, err := s.Client.GetConversation(ctx, &chatpb.ConversationPostRequest{
// 		XUserId:         _userId,
// 		ConversationIds: conversationIds,
// 	})
// 	if err != nil {
// 		err = xerrors.Errorf("%w", err)
// 		lib.RecordError(span, err, ctx)
// 		return nil, err
// 	}
// 	return res.Data, nil
// }
