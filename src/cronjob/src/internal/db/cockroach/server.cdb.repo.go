package cockroach

import (
	"time"

	"github.com/aqaurius6666/cronjob/src/internal/db/payment"
	"github.com/aqaurius6666/cronjob/src/internal/db/transaction"
	"github.com/aqaurius6666/go-utils/database/cockroach"
	"github.com/google/wire"
)

var (
	timeout = 10 * time.Second
)

// var (
// 	_ db.ServerRepo = (*ServerCDBRepo)(nil)
// )
var CDBRepoSet = wire.NewSet(wire.Struct(new(ServerCDBRepo), "*"), InterfacesProvider, wire.Struct(new(cockroach.CDBRepository), "*"))

func InterfacesProvider() cockroach.DBInterfaces {
	return cockroach.DBInterfaces{
		payment.Payment{},
		transaction.Transaction{},
	}
}

type ServerCDBRepo struct {
	cockroach.CDBRepository
}
