package model

import (
	"context"

	"github.com/aqaurius6666/cronjob/src/internal/db"
	"github.com/aqaurius6666/cronjob/src/services/chatservice"
	"github.com/aqaurius6666/cronjob/src/services/mailservice"
	"github.com/aqaurius6666/cronjob/src/services/payment"
	"github.com/aqaurius6666/cronjob/src/services/redis"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

type Server interface {
	TransactionModel
	NotificationModel
	ConversationModel
}

type ServerModel struct {
	Ctx     context.Context
	Logger  *logrus.Logger
	Repo    db.ServerRepo
	Payment payment.Interface
	Mail    mailservice.Service
	Redis   redis.Redis
	Chat    chatservice.Service
}

var ServerModelSet = wire.NewSet(wire.Bind(new(Server), new(*ServerModel)), wire.Struct(new(ServerModel), "*"))
