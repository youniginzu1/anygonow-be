package business

import (
	"github.com/aqaurius6666/apiservice/src/internal/db/contact"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/ubgo/gormuuid"
)

type Business struct {
	database.BaseModel
	Name           *string            `gorm:"type:varchar(64)"`
	Phone          *string            `gorm:"type:varchar(64)"`
	LogoUrl        *string            `gorm:"type:text"`
	BannerUrl      *string            `gorm:"type:text"`
	Mail           *string            `gorm:"type:varchar(64)"`
	ContactId      uuid.UUID          `gorm:"type:uuid"`
	Contact        *contact.Contact   `gorm:"foreginKey:ContactId"`
	Website        *string            `gorm:"type:text"`
	Description    *string            `gorm:"type:varchar(1024)"`
	Services       gormuuid.UUIDArray `gorm:"type:uuid[]"`
	Zipcodes       pq.StringArray     `gorm:"type:varchar(16)[]"`
	Status         *int32             `gorm:"type:int8;default:0"`
	InvitationCode *string            `gorm:"type:varchar(64)"`
	RefCode        *string            `gorm:"type:varchar(64)"`
	FreeContact    *int32             `gorm:"type:int8;default:0"`
	RefStatus      *int32             `gorm:"type:int8;default:0"`
	Rate           *float32           `gorm:"-:migration;->"`
	Review         *int32             `gorm:"-:migration;->"`
	Request        *int32             `gorm:"-:migration;->"`
	Zipcode        *string            `gorm:"-:migration;->"`
	ServiceName    pq.StringArray     `gorm:"type:varchar(64)[];-:migration;->"`
	ServiceId      pq.StringArray     `gorm:"type:varchar(64)[];-:migration;->"`
	StartDate      *int64             `gorm:"-:migration;->"`
	CountZipcodes  *int64             `gorm:"-:migration;->"`
}

type Search struct {
	database.DefaultSearchModel
	Business
	Zipcode      *string
	CategoryId   uuid.UUID
	Query        c.SORT_QUERY
	BothId       []string
	ValidateMail *string
}
