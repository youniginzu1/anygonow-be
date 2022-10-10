package otp

import (
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
)

type Otp struct {
	database.BaseModel
	Type       *int
	Code       *string
	Status     *int
	ExpireTime int64
	Metadata   *string
	UserId     uuid.UUID `gorm:"type:uuid"`
	// User       *user.User `gorm:"foreignKey:UserId"`
	Mail     *string
	Commited *bool `gorm:"default:false"`
}

type Search struct {
	database.DefaultSearchModel
	Otp
}
