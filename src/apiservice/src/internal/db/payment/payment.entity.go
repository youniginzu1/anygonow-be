package payment

import (
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
)

type Payment struct {
	database.BaseModel
	BusinessId      uuid.UUID `gorm:"type:uuid"`
	CustomerId      *string   `gorm:"type:varchar(32)"`
	PaymentMethodId *string   `gorm:"type:varchar(64)"`
}
type Search struct {
	database.DefaultSearchModel
	Payment
}
