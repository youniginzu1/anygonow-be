package notification

import (
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
)

type Notification struct {
	database.BaseModel
	UserId uuid.UUID `gorm:"type:uuid;unique"`
	Seen   *bool     `gorm:"type:bool;default:false"`
}

type Search struct {
	database.DefaultSearchModel
	Notification
}
