package advertise_order

import (
	"github.com/aqaurius6666/apiservice/src/internal/db/advertise_package"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
)

type AdvertiseOrder struct {
	database.BaseModel
	BusinessId             uuid.UUID                          `gorm:"type:uuid"`
	ServiceId              uuid.UUID                          `gorm:"type:uuid"`
	StartDate              *int64                             `gorm:"type:bigint"`
	EndDate                *int64                             `gorm:"type:bigint"`
	AdvertisePackageId     uuid.UUID                          `gorm:"type:uuid"`
	AdvertiseTransactionId uuid.UUID                          `gorm:"type:uuid"`
	AdvertisePackage       advertise_package.AdvertisePackage `gorm:"foreignKey:AdvertisePackageId"`
	Name                   *string                            `gorm:"-:migration;->"`
	Price                  *float64                           `gorm:"-:migration;->"`
	BannerUrl              *string                            `gorm:"-:migration;->"`
	Description            *string                            `gorm:"-:migration;->"`
	Zipcode                *string                            `gorm:"-:migration;->"`
	CategoryName           *string                            `gorm:"-:migration;->"`
}

type Search struct {
	database.DefaultSearchModel
	AdvertiseOrder
	CategoryId uuid.UUID 
	Zipcode *string
}
