//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/aqaurius6666/authservice/src/internal/api"
	"github.com/aqaurius6666/authservice/src/internal/db"
	"github.com/aqaurius6666/authservice/src/internal/model"
	"github.com/aqaurius6666/authservice/src/services/mailservice"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

type Server struct {
	ApiServer *api.ApiServer
	MainRepo  db.ServerRepo
}

type ServerOptions struct {
	DBDsn    db.DBDsn
	MailAddr mailservice.MailServiceAddr
}

func InitMainServer(ctx context.Context, logger *logrus.Logger, opts ServerOptions) (*Server, error) {

	wire.Build(
		wire.FieldsOf(&opts, "DBDsn", "MailAddr"),
		db.ServerRepoSet,
		api.ApiServerSet,
		model.ServerModelSet,
		mailservice.MailServiceSet,
		wire.Struct(new(Server), "*"),
	)
	return &Server{}, nil
}
