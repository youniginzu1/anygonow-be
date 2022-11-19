package conversation

import (
	"github.com/aqaurius6666/chatservice/src/internal/db/phone_pool"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/ubgo/gormuuid"
)

type Conversation struct {
	database.BaseModel
	Members            gormuuid.UUIDArray    `gorm:"type:uuid[]"`
	PhoneNumberMembers pq.StringArray        `gorm:"type:varchar(16)[]"`
	OrderId            uuid.UUID             `gorm:"type:uuid"`
	Status             *int32                `gorm:"type:int8;default:0"`
	LastChat           *int64                `gorm:"-:migration;->"`
	MemberNames        pq.StringArray        `gorm:"type:varchar(64)[];-:migration;->"`
	PhonePoolId        uuid.UUID             `gorm:"type:uuid"`
	PhonePool          *phone_pool.PhonePool `gorm:"foreignKey:PhonePoolId"`
	ServiceId          uuid.UUID             `gorm:"type:uuid"`
}

type Search struct {
	database.DefaultSearchModel
	Conversation
	MemberId        uuid.UUID
	ConversationIds gormuuid.UUIDArray
	MemberPhone     *string
}
