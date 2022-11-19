package order

import (
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
	"github.com/ubgo/gormuuid"
)

type Order struct {
	database.BaseModel
	CustomerId      uuid.UUID `gorm:"type:uuid"`
	BusinessId      uuid.UUID `gorm:"type:uuid"`
	ConversationId  uuid.UUID `gorm:"type:uuid"`
	ServiceId       uuid.UUID `gorm:"type:uuid"`
	StartDate       *int64    `gorm:"type:bigint"`
	EndDate         *int64    `gorm:"type:bigint"`
	Status          *int32    `gorm:"type:int8;default:0"`
	CustomerPhone   *string   `gorm:"type:varchar(16)"`
	CustomerName    *string   `gorm:"type:varchar(128)"`
	CustomerZipcode *string   `gorm:"type:varchar(16)"`
	CustomerMessage *string   `gorm:"type:varchar(256)"`
	IsReviewed      *bool     `gorm:"type:bool;default:false"`
	ServiceName     *string   `gorm:"-:migration;->"`
	NumberOrders    *int64    `gorm:"-:migration;->"`
	ServiceAvatar   *string   `gorm:"-:migration;->"`
	CustomerAvatar  *string   `gorm:"-:migration;->"`
	Fee             *float32  `gorm:"-:migration;->"`
	CategoryName    *string   `gorm:"-:migration;->"`
	BusinessName    *string   `gorm:"-:migration;->"`
	BusinessLogo    *string   `gorm:"-:migration;->"`
	BusinessBanner  *string   `gorm:"-:migration;->"`
	CategoryId      *string   `gorm:"-:migration;->"`
	HandymanMail    *string   `gorm:"-:migration;->"`
	CustomerMail    *string   `gorm:"-:migration;->"`
}
type Search struct {
	database.DefaultSearchModel
	Order
	UserId     uuid.UUID
	CategoryId uuid.UUID
	Query      *int32
	OIds       gormuuid.UUIDArray
}
