package mailservice

import (
	"context"
	"time"

	"github.com/aqaurius6666/authservice/src/pb/mailpb"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	_       Service = ServiceGRPC{}
	timeout         = 5 * time.Second
)

type MailServiceAddr string

type ServiceGRPC struct {
	Ctx    context.Context
	Client mailpb.MailServiceClient
}

func ConnectClient(ctx context.Context, addr MailServiceAddr) (mailpb.MailServiceClient, error) {
	nctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	conn, err := grpc.DialContext(nctx, string(addr), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}
	return mailpb.NewMailServiceClient(conn), nil
}

func (s ServiceGRPC) SendMail(ctx context.Context, to string, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	_, err := s.Client.SendMail(ctx, &mailpb.SendMailRequest{
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

func (s ServiceGRPC) SendMails(ctx context.Context, to []string, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	_, err := s.Client.SendMails(ctx, &mailpb.SendMailsRequest{
		Tos: to,
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
