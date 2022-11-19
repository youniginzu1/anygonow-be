package contact

import (
	"github.com/aqaurius6666/apiservice/src/internal/db/state"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
)

type Contact struct {
	database.BaseModel
	Zipcode   *string      `gorm:"type:varchar(16)"`
	StateId   uuid.UUID    `gorm:"type:uuid"`
	State     *state.State `gorm:"foreignKey:StateId"`
	Address1  *string      `gorm:"type:varchar(256)"`
	Address2  *string      `gorm:"type:varchar(256)"`
	City      *string      `gorm:"type:varchar(64)"`
}

type Search struct {
	database.DefaultSearchModel
	Contact
}
