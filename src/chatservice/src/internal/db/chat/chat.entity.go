package chat

import (
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
)

type Chat struct {
	database.BaseModel
	SenderId       uuid.UUID `gorm:"type:uuid"`
	ConversationId uuid.UUID `gorm:"type:uuid"`
	Payload        *string   `gorm:"type:text"`
	Seen           *bool     `gorm:"default:false"`
}

type Search struct {
	database.DefaultSearchModel
	Chat
	Timestamp  *int64
	Min        *int32
	ReceiverId uuid.UUID
}
