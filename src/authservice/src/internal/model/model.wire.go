//go:build wireinject
// +build wireinject

package model

import (
	"context"

	"github.com/aqaurius6666/authservice/src/internal/db"
	"github.com/aqaurius6666/authservice/src/services/mailservice"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

type ModelMockOptions struct {
	DBDsn    db.DBDsn
	MailAddr mailservice.MailServiceAddr
}

func ModelMock(ctx context.Context, logger *logrus.Logger, opts ModelMockOptions) (Server, error) {

	wire.Build(
		wire.FieldsOf(&opts, "DBDsn", "MailAddr"),
		db.ServerRepoSet,
		mailservice.MailServiceSet,
		wire.Struct(new(ServerModel), "*"),
		wire.Bind(new(Server), new(*ServerModel)),
	)
	return nil, nil
}
