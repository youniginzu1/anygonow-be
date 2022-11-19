package model

import (
	"context"

	"github.com/aqaurius6666/apiservice/src/internal/db"
	"github.com/aqaurius6666/apiservice/src/pb/authpb"
	"github.com/aqaurius6666/apiservice/src/services/authservice"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SetUpRepo() (db.ServerRepo, error) {
	return nil, nil
}
func SetUpModel() (Server, error) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	repo, err := SetUpRepo()
	if err != nil {
		return nil, err
	}

	model := &ServerModel{
		Ctx:    context.Background(),
		Logger: logrus.New(),
		Repo:   repo,
		Auth: authservice.ServiceGRPC{
			Ctx:    context.Background(),
			Client: authpb.NewAuthServiceClient(conn),
		},
	}
	return model, nil
}
