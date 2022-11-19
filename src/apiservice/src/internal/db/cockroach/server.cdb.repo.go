package cockroach

import (
	"time"

	"github.com/aqaurius6666/apiservice/src/internal/db/advertise_order"
	"github.com/aqaurius6666/apiservice/src/internal/db/advertise_package"
	"github.com/aqaurius6666/apiservice/src/internal/db/advertise_transaction"
	"github.com/aqaurius6666/apiservice/src/internal/db/business"
	"github.com/aqaurius6666/apiservice/src/internal/db/category"
	"github.com/aqaurius6666/apiservice/src/internal/db/contact"
	"github.com/aqaurius6666/apiservice/src/internal/db/feedback"
	"github.com/aqaurius6666/apiservice/src/internal/db/group"
	"github.com/aqaurius6666/apiservice/src/internal/db/order"
	"github.com/aqaurius6666/apiservice/src/internal/db/payment"
	"github.com/aqaurius6666/apiservice/src/internal/db/service"
	"github.com/aqaurius6666/apiservice/src/internal/db/state"
	"github.com/aqaurius6666/apiservice/src/internal/db/transaction"
	"github.com/aqaurius6666/apiservice/src/internal/db/user"
	"github.com/aqaurius6666/go-utils/database/cockroach"
	"github.com/google/wire"
)

var (
	timeout = 2 * time.Second
)

// var (
// 	_ db.ServerRepo = (*ServerCDBRepo)(nil)
// )
var CDBRepoSet = wire.NewSet(wire.Struct(new(ServerCDBRepo), "*"), InterfacesProvider, wire.Struct(new(cockroach.CDBRepository), "*"))

func InterfacesProvider() cockroach.DBInterfaces {
	return cockroach.DBInterfaces{
		advertise_order.AdvertiseOrder{},
		advertise_package.AdvertisePackage{},
		advertise_transaction.AdvertiseTransaction{},
		business.Business{},
		category.Category{},
		contact.Contact{},
		feedback.Feedback{},
		order.Order{},
		service.Service{},
		state.State{},
		user.User{},
		payment.Payment{},
		group.Group{},
		transaction.Transaction{},
	}
}

type ServerCDBRepo struct {
	cockroach.CDBRepository
}
