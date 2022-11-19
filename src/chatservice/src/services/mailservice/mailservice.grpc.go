package mailservice

import (
	"context"
	"time"

	"github.com/aqaurius6666/chatservice/src/internal/lib"
	"github.com/aqaurius6666/chatservice/src/internal/var/c"
	"github.com/aqaurius6666/chatservice/src/pb/mailpb"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type MailserviceAddr string

var (
	_       Service = (*ServiceGRPC)(nil)
	timeout         = 5 * time.Second
)

type ServiceGRPC struct {
	Ctx    context.Context
	Client mailpb.MailServiceClient
}

func ConnectClient(ctx context.Context, addr MailserviceAddr) (mailpb.MailServiceClient, error) {
	nctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	conn, err := grpc.DialContext(nctx, string(addr), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithChainUnaryInterceptor(
		otelgrpc.UnaryClientInterceptor(),
	))
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}
	return mailpb.NewMailServiceClient(conn), nil
}
func (s *ServiceGRPC) SubscribeNotification(ctx context.Context, userId string, deviceId string) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.SubscribeNotification))
	defer span.End()

	nctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	_, err := s.Client.SubscribeNotification(nctx, &mailpb.SubscribeNotificationRequest{
		UserId:   userId,
		DeviceId: deviceId,
	})
	if err != nil {
		stt, ok := status.FromError(err)
		if ok {
			err = xerrors.New(stt.Message())
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)

		return err
	}
	return nil
}

func (s *ServiceGRPC) UnsubscribeNotification(ctx context.Context, userId string, deviceId string) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UnsubscribeNotification))
	defer span.End()

	nctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	_, err := s.Client.UnsubscribeNotification(nctx, &mailpb.UnsubscribeNotificationRequest{
		UserId:   userId,
		DeviceId: deviceId,
	})
	if err != nil {
		stt, ok := status.FromError(err)
		if ok {
			err = xerrors.New(stt.Message())
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)

		return err
	}
	return nil
}
func (s *ServiceGRPC) SendNotification(ctx context.Context, to string, title, body string, message string) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.SendNotification))
	defer span.End()
	nctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	_, err := s.Client.SendNotification(nctx, &mailpb.SendNotificationRequest{
		To:      to,
		Body:    body,
		Message: message,
		Title:   title,
	})
	if err != nil {
		stt, ok := status.FromError(err)
		if ok {
			err = xerrors.New(stt.Message())
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)

		return err
	}
	return nil
}
func (s *ServiceGRPC) SendMail(to string, msg []byte) error {
	_, err := s.Client.SendMail(s.Ctx, &mailpb.SendMailRequest{
		To:  to,
		Msg: msg,
	})

	if err != nil {
		stt, ok := status.FromError(err)
		if !ok {
			return xerrors.Errorf("%w", err)
		}
		return xerrors.Errorf("%w", xerrors.Errorf("%w", stt.Err()))
	}
	return nil
}
