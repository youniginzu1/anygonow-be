package api

import (
	"context"

	"github.com/aqaurius6666/mailservice/src/pb/mailpb"
	"golang.org/x/xerrors"
)

func (s *ApiServer) SendMail(ctx context.Context, req *mailpb.SendMailRequest) (*mailpb.SendMailResponse, error) {
	s.Logger.Infof("Send mail to %s", req.To)
	err := s.MailService.SendMail(req.To, req.Msg)
	if err != nil {
		s.Logger.Info(string(req.Msg))
		panic(xerrors.Errorf("%w", err))
	}
	return &mailpb.SendMailResponse{}, nil
}

func (s *ApiServer) SendMails(ctx context.Context, req *mailpb.SendMailsRequest) (*mailpb.SendMailsResponse, error) {
	err := s.MailService.SendMails(req.Tos, req.Msg)
	if err != nil {
		s.Logger.Info(string(req.Msg))
		panic(xerrors.Errorf("%w", err))
	}
	return &mailpb.SendMailsResponse{}, nil
}
