//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/aqaurius6666/cronjob/src/internal/cronjob"
	"github.com/aqaurius6666/cronjob/src/internal/db"
	"github.com/aqaurius6666/cronjob/src/internal/model"
	"github.com/aqaurius6666/cronjob/src/services/chatservice"
	"github.com/aqaurius6666/cronjob/src/services/mailservice"
	"github.com/aqaurius6666/cronjob/src/services/payment"
	"github.com/aqaurius6666/cronjob/src/services/redis"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Jobs     []cronjob.Cronjob
	MainRepo db.ServerRepo
}

type ServerOptions struct {
	DBDsn                db.DBDsn
	Key                  payment.STRIPE_API_KEY
	QuantityInterval     cronjob.QUANTITY_INTERVAL
	ChatQuantityInterval cronjob.CHAT_QUANTITY_INTERVAL
	UnitInterval         cronjob.UNIT_INTERVAL
	PaymentDay           cronjob.PAYMENT_DAY
	MailserviceAddr      mailservice.MailserviceAddr
	RedisUri             redis.REDIS_URI
	RedisUser            redis.REDIS_USER
	RedisPass            redis.REDIS_PASS
	ChatserviceAddr      chatservice.ChatserviceAddr
}

func InitMainServer(ctx context.Context, logger *logrus.Logger, opts ServerOptions) (*Server, error) {

	wire.Build(
		wire.FieldsOf(&opts, "DBDsn", "Key", "QuantityInterval", "UnitInterval", "PaymentDay", "MailserviceAddr", "ChatQuantityInterval", "RedisUri", "RedisUser", "RedisPass", "ChatserviceAddr"),
		db.ServerRepoSet,
		model.ServerModelSet,
		cronjob.JobSet,
		payment.Set,
		mailservice.Set,
		redis.Set,
		chatservice.Set,
		wire.Struct(new(Server), "*"),
	)
	return &Server{}, nil
}
