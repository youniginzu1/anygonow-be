package model

import (
	"context"

	"github.com/aqaurius6666/apiservice/src/internal/db"
	"github.com/aqaurius6666/apiservice/src/services/authservice"
	"github.com/aqaurius6666/apiservice/src/services/chatservice"
	"github.com/aqaurius6666/apiservice/src/services/mailservice"
	"github.com/aqaurius6666/apiservice/src/services/payment"
	"github.com/aqaurius6666/apiservice/src/services/s3"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

type Server interface {
	UserModel
	BusinessModel
	ContactModel
	AdminModel
	S3Model
	FeedbackModel
	ServiceModel
	OrderModel
	CategoryModel
	PaymentModel
	ChatModel
	GroupModel
	AdvertisePackageModel
	AdvertiseOrderModel
	TransactionModel
	NotificationModel
}

type ServerModel struct {
	Ctx     context.Context
	Logger  *logrus.Logger
	Repo    db.ServerRepo
	Auth    authservice.Service
	S3      s3.Service
	Payment payment.Interface
	Chat    chatservice.Service
	Mail    mailservice.Service
}

var ServerModelSet = wire.NewSet(wire.Bind(new(Server), new(*ServerModel)), wire.Struct(new(ServerModel), "*"))
