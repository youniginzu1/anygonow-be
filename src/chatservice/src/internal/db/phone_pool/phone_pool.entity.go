package phone_pool

import (
	"github.com/aqaurius6666/go-utils/database"
)

type PhonePool struct {
	database.BaseModel
	PhoneNumber *string
	Status      *int32 `gorm:"type:int8;default:0"`
	Sid         *string
}

type Search struct {
	database.DefaultSearchModel
	PhonePool
}
