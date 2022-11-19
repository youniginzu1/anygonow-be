package middleware

import (
	"github.com/aqaurius6666/chatservice/src/services/authservice"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

var Set = wire.NewSet(wire.Bind(new(Middleware), new(*MiddlewareV1)), wire.Struct(new(MiddlewareV1), "*"))

type Middleware interface {
	Auth
	Logger
	Body
	// Role
}

type MiddlewareV1 struct {
	Logger *logrus.Logger
	Auth   authservice.Service
}
