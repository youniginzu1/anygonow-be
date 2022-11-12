//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/aqaurius6666/mailservice/src/internal/api"
	"github.com/aqaurius6666/mailservice/src/internal/db"
	"github.com/aqaurius6666/mailservice/src/internal/mail"
	"github.com/aqaurius6666/mailservice/src/internal/model"
	"github.com/aqaurius6666/mailservice/src/services/fcm"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

type Server struct {
	ApiServer *api.ApiServer
	Repo      db.ServerRepo
}

type ServerOptions struct {
	Username    mail.MailUsername
	Password    mail.MailPassword
	Host        mail.SMTPHost
	Port        mail.SMTPPort
	Sender      mail.SMTPSender
	DBDsn       db.DBDsn
	FirebaseKey fcm.FB_PRIVATE_KEY
}

func InitMainServer(ctx context.Context, logger *logrus.Logger, opts ServerOptions) (*Server, error) {

	wire.Build(
		wire.FieldsOf(&opts, "Username", "Password", "Host", "Port", "Sender", "DBDsn", "FirebaseKey"),
		api.ApiServerSet,
		mail.ServiceSet,
		model.ServerModelSet,
		fcm.Set,
		db.ServerRepoSet,
		wire.Struct(new(Server), "*"),
	)
	return &Server{}, nil
}
