package model

import (
	"context"

	"github.com/aqaurius6666/chatservice/src/internal/db"
	"github.com/aqaurius6666/chatservice/src/services/mailservice"
	"github.com/aqaurius6666/chatservice/src/services/redis"
	"github.com/aqaurius6666/chatservice/src/services/twilloclient"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

var (
	_ Server = (*ServerModel)(nil)
)

type Server interface {
	HTTP
	Conversation
	Notification
}

type ServerModel struct {
	Ctx    context.Context
	Logger *logrus.Logger
	Repo   db.ServerRepo
	Twillo twilloclient.Twilio
	Redis  redis.Redis
	Mail   mailservice.Service
}

var ServerModelSet = wire.NewSet(wire.Bind(new(Server), new(*ServerModel)), wire.Struct(new(ServerModel), "*"))
