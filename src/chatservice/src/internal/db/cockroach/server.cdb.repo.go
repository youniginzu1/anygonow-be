package cockroach

import (
	"time"

	"github.com/aqaurius6666/chatservice/src/internal/db/chat"
	"github.com/aqaurius6666/chatservice/src/internal/db/conversation"
	"github.com/aqaurius6666/chatservice/src/internal/db/notification"
	"github.com/aqaurius6666/chatservice/src/internal/db/phone_pool"
	"github.com/aqaurius6666/go-utils/database/cockroach"
	"github.com/google/wire"
)

// var (
// 	_ db.ServerRepo = (*ServerCDBRepo)(nil)
// )
var (
	timeout = 2 * time.Second
)
var CDBRepoSet = wire.NewSet(wire.Struct(new(ServerCDBRepo), "*"), InterfacesProvider, wire.Struct(new(cockroach.CDBRepository), "*"))

func InterfacesProvider() cockroach.DBInterfaces {
	return cockroach.DBInterfaces{
		chat.Chat{},
		conversation.Conversation{},
		notification.Notification{},
		phone_pool.PhonePool{},
	}
}

type ServerCDBRepo struct {
	cockroach.CDBRepository
}
