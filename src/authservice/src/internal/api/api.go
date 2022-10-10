package api

import (
	"github.com/aqaurius6666/authservice/src/internal/db"
	"github.com/aqaurius6666/authservice/src/internal/model"
	"github.com/aqaurius6666/authservice/src/pb/authpb"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

var ApiServerSet = wire.NewSet(wire.Struct(new(ApiServer), "*"))

var (
	_ authpb.AuthServiceServer = (*ApiServer)(nil)
)

type ApiServer struct {
	authpb.UnimplementedAuthServiceServer `wire:"-"`
	Model                                 model.Server
	Logger                                *logrus.Logger
	Repo                                  db.ServerRepo
}
