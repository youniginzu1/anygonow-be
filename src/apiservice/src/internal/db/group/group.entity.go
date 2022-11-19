package group

import (
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/ubgo/gormuuid"
)

type Group struct {
	database.BaseModel
	Name        *string            `gorm:"type:varchar(64)"`
	Fee         *float32           `gorm:"type:float8;default:0"`
	CategoryIds gormuuid.UUIDArray `gorm:"type:uuid[]"`
	ServiceName pq.StringArray     `gorm:"type:varchar(64)[];-:migration;->"`
	ServiceId   pq.StringArray     `gorm:"type:varchar(64)[];-:migration;->"`
}

type Search struct {
	database.DefaultSearchModel
	Group
	CategoryId uuid.UUID
}
