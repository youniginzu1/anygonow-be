package db

import (
	"github.com/aqaurius6666/chatservice/src/internal/db/chat"
	"github.com/aqaurius6666/chatservice/src/internal/db/conversation"
	"github.com/aqaurius6666/chatservice/src/internal/db/notification"
	"github.com/aqaurius6666/chatservice/src/internal/db/phone_pool"
	"github.com/aqaurius6666/go-utils/database"
)

type DBDsn string

type ServerRepo interface {
	database.CommonRepository
	conversation.ConversationRepo
	chat.ChatRepo
	notification.NotificationRepo
	phone_pool.PhonePoolRepo
}
