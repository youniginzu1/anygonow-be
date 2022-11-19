package db

import (
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
	"github.com/aqaurius6666/go-utils/database"
)

type DBDsn string
type ServerRepo interface {
	database.CommonRepository
	user.UserRepo
	business.BusinessRepo
	contact.ContactRepo
	service.ServiceRepo
	state.StateRepo
	advertise_package.AdvertisePackageRepo
	feedback.FeedbackRepo
	category.CategoryRepo
	order.OrderRepo
	payment.PaymentRepo
	group.GroupRepo
	transaction.TransactionRepo
	advertise_order.AdvertiseOrderRepo
	advertise_transaction.AdvertiseTransactionRepo
}
