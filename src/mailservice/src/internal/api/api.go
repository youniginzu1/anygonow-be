package api

import (
	"github.com/aqaurius6666/mailservice/src/internal/mail"
	"github.com/aqaurius6666/mailservice/src/internal/model"
	"github.com/aqaurius6666/mailservice/src/pb/mailpb"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

var ApiServerSet = wire.NewSet(wire.Struct(new(ApiServer), "*"))

var (
	_ mailpb.MailServiceServer = (*ApiServer)(nil)
)

type ApiServer struct {
	mailpb.UnimplementedMailServiceServer `wire:"-"`
	MailService                           mail.Service
	Logger                                *logrus.Logger
	Model                                 model.Server
}
