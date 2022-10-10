package model

import (
	"context"

	"github.com/aqaurius6666/authservice/src/internal/db"
	"github.com/aqaurius6666/authservice/src/services/mailservice"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

var (
	_ Server = (*ServerModel)(nil)
)

type Server interface {
	UserModel
	OtpModel
	RoleModel
}

type ServerModel struct {
	Ctx    context.Context
	Logger *logrus.Logger
	Repo   db.ServerRepo
	Mail   mailservice.Service
}

var ServerModelSet = wire.NewSet(wire.Bind(new(Server), new(*ServerModel)), wire.Struct(new(ServerModel), "*"))
