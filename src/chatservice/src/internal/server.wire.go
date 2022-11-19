//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/aqaurius6666/chatservice/src/internal/api"
	"github.com/aqaurius6666/chatservice/src/internal/db"
	"github.com/aqaurius6666/chatservice/src/internal/model"
	"github.com/aqaurius6666/chatservice/src/services/authservice"
	"github.com/aqaurius6666/chatservice/src/services/mailservice"
	"github.com/aqaurius6666/chatservice/src/services/redis"
	"github.com/aqaurius6666/chatservice/src/services/twilloclient"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

type Server struct {
	ApiServer *api.ApiServer
	MainRepo  db.ServerRepo
}

type ServerOptions struct {
	DBDsn             db.DBDsn
	AuthserviceAddr   authservice.AuthserviceAddr
	MailserviceAddr   mailservice.MailserviceAddr
	RedisUri          redis.REDIS_URI
	RedisUser         redis.REDIS_USER
	RedisPass         redis.REDIS_PASS
	TwilloCallbackUrl twilloclient.TWILLO_SMS_CALLBACK_URL
}

func InitMainServer(ctx context.Context, logger *logrus.Logger, opts ServerOptions) (*Server, error) {
	wire.Build(
		wire.FieldsOf(&opts, "TwilloCallbackUrl", "DBDsn", "AuthserviceAddr", "RedisUser", "RedisPass", "RedisUri", "MailserviceAddr"),
		gin.New,
		db.ServerRepoSet,
		api.ApiServerSet,
		authservice.Set,
		model.ServerModelSet,
		twilloclient.Set,
		redis.Set,
		mailservice.Set,
		wire.Struct(new(Server), "*"),
	)
	return &Server{}, nil
}
