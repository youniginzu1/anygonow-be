package model

import (
	"context"

	"github.com/aqaurius6666/mailservice/src/internal/db"
	"github.com/aqaurius6666/mailservice/src/services/fcm"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

var (
	_ Server = (*ServerModel)(nil)
)

type Server interface {
	NotificationModel
}

type ServerModel struct {
	Ctx    context.Context
	Logger *logrus.Logger
	Repo   db.ServerRepo
	Fcm    fcm.Service
}

var ServerModelSet = wire.NewSet(wire.Bind(new(Server), new(*ServerModel)), wire.Struct(new(ServerModel), "*"))
