//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/aqaurius6666/apiservice/src/internal/api"
	"github.com/aqaurius6666/apiservice/src/internal/db"
	"github.com/aqaurius6666/apiservice/src/internal/model"
	"github.com/aqaurius6666/apiservice/src/services/authservice"
	"github.com/aqaurius6666/apiservice/src/services/chatservice"
	"github.com/aqaurius6666/apiservice/src/services/mailservice"
	"github.com/aqaurius6666/apiservice/src/services/payment"
	"github.com/aqaurius6666/apiservice/src/services/s3"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

type Server struct {
	ApiServer *api.ApiServer
	MainRepo  db.ServerRepo
}

type ServerOptions struct {
	DBDsn           db.DBDsn
	ChatserviceAddr chatservice.ChatserviceAddr
	AuthserviceAddr authservice.AuthServiceAddr
	MailserviceAddr mailservice.MailserviceAddr
	Bucket          s3.BucketName
	Key             payment.STRIPE_API_KEY
	SignKey         api.STRIPE_SIGNATURE_KEY
}

func InitMainServer(ctx context.Context, logger *logrus.Logger, opts ServerOptions) (*Server, error) {

	wire.Build(
		wire.FieldsOf(&opts, "DBDsn", "AuthserviceAddr", "Bucket", "Key", "ChatserviceAddr", "SignKey", "MailserviceAddr"),
		db.ServerRepoSet,
		api.ApiServerSet,
		model.ServerModelSet,
		authservice.Set,
		chatservice.Set,
		payment.Set,
		s3.Set,
		mailservice.Set,
		wire.Struct(new(Server), "*"),
	)
	return &Server{}, nil
}
