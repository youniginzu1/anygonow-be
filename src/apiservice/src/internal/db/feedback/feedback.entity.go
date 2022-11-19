package feedback

import (
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
)

type Feedback struct {
	database.BaseModel
	Rate         *float32  `gorm:"type:float8"`
	Comment      *string   `gorm:"type:varchar(512)"`
	ServiceId    uuid.UUID `gorm:"type:uuid"`
	BusinessId   uuid.UUID `gorm:"type:uuid"`
	Review       *int32    `gorm:"-:migration;->"`
	UserId       uuid.UUID `gorm:"type:uuid"`
	OrderId      uuid.UUID `gorm:"type:uuid"`
	CustomerName *string   `gorm:"-:migration;->"`
	AvatarUrl    *string   `gorm:"-:migration;->"`
	ServiceName  *string   `gorm:"-:migration;->"`
	Request      *int32    `gorm:"-:migration;->"`
}

type Search struct {
	database.DefaultSearchModel
	Feedback
}
