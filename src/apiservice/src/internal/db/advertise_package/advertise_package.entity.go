package advertise_package

import (
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/ubgo/gormuuid"
)

type AdvertisePackage struct {
	database.BaseModel
	Name        *string
	Categories  gormuuid.UUIDArray `gorm:"type:uuid[]"`
	Price       *float64
	BannerUrl   *string
	Description *string
	ServiceName pq.StringArray `gorm:"type:varchar(64)[];-:migration;->"`
	ServiceId   pq.StringArray `gorm:"type:varchar(64)[];-:migration;->"`
	Zipcodes    pq.StringArray `gorm:"type:varchar(16)[];-:migration;->"`
}

type Search struct {
	database.DefaultSearchModel
	AdvertisePackage
	ServiceName string
	BusinessId  uuid.UUID
}
