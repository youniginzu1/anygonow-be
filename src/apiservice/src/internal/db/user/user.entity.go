package user

import (
	"github.com/aqaurius6666/apiservice/src/internal/db/contact"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
)

type User struct {
	database.BaseModel
	Mail      *string          `gorm:"varchar(64)"`
	Phone     *string          `gorm:"type:varchar(64)"`
	ContactId uuid.UUID        `gorm:"type:uuid"`
	Contact   *contact.Contact `gorm:"foreignKey:ContactId"`
	AvatarUrl *string          `gorm:"type:text"`
	FirstName *string          `gorm:"type:varchar(128)"`
	LastName  *string          `gorm:"type:varchar(128)"`
	Status    *int32           `gorm:"type:int8;default:0"`
}

type Search struct {
	database.DefaultSearchModel
	User
	Query *string
}
