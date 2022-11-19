package transaction

import (
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
)

type Transaction struct {
	database.BaseModel
	OrderId          uuid.UUID `gorm:"type:uuid"`
	BusinessId       uuid.UUID `gorm:"type:uuid"`
	Fee              *int64    `gorm:"type:int8;default:0"`
	IsFree           *bool     `gorm:"type:bool;default:false"`
	IsPaid           *bool     `gorm:"type:bool;default:false"`
	PaymentIntentId  *string   `gorm:"type:varchar(128)"`
	StartDate        *int64    `gorm:"-:migration;->"`
	EndDate          *int64    `gorm:"-:migration;->"`
	ServiceName      *string   `gorm:"-:migration;->"`
	CustomerZipcode  *string   `gorm:"-:migration;->"`
	Status           *int32    `gorm:"-:migration;->"`
	CustomerId       uuid.UUID `gorm:"-:migration;->"`
	CustomerAvatar   *string   `gorm:"-:migration;->"`
	StripeCustomerId *string   `gorm:"-:migration;->"`
	PaymentMethodId  *string   `gorm:"-:migration;->"`
}

type Search struct {
	database.DefaultSearchModel
	Transaction
	Query *int32
}
