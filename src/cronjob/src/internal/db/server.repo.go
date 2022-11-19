package db

import (
	"github.com/aqaurius6666/cronjob/src/internal/db/payment"
	"github.com/aqaurius6666/cronjob/src/internal/db/transaction"
	"github.com/aqaurius6666/go-utils/database"
)

type DBDsn string

type ServerRepo interface {
	database.CommonRepository
	payment.PaymentRepo
	transaction.TransactionRepo
}
