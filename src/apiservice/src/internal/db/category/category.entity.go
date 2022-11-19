package category

import (
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
	"github.com/ubgo/gormuuid"
)

type Category struct {
	database.BaseModel
	Name          *string   `gorm:"type:varchar(64)"`
	GroupId       uuid.UUID `gorm:"type:uuid"`
	ImageUrl      *string   `gorm:"type:varchar(512)"`
	Fee           *float32  `gorm:"-:migration;->"`
	TotalProvider *int64    `gorm:"-:migration;->"`
}

type Search struct {
	database.DefaultSearchModel
	Category
	CategoryIds gormuuid.UUIDArray
	Query       c.QUERY_CATEGORY_ADMIN
}
