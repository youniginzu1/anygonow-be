//go:build wireinject
// +build wireinject

package model

import (
	"context"

	"github.com/aqaurius6666/chatservice/src/internal/db"
	"github.com/aqaurius6666/chatservice/src/services/mailservice"
	"github.com/aqaurius6666/chatservice/src/services/redis"
	"github.com/aqaurius6666/chatservice/src/services/twilloclient"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"testing"
)

func NewMockModel(ctx context.Context, logger *logrus.Logger, dsn db.DBDsn, t *testing.T) (Server, error) {
	wire.Build(
		ServerModelSet,
		db.NewMockRepo,
		wire.Struct(new(mailservice.ServiceGRPC)),
		wire.Bind(new(mailservice.Service), new(*mailservice.ServiceGRPC)),
		wire.Struct(new(redis.RedisImpl)),
		wire.Bind(new(redis.Redis), new(*redis.RedisImpl)),
		twilloclient.MockSet,
	)
	return nil, nil
}
