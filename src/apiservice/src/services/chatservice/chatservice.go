package chatservice

import (
	"context"

	"github.com/aqaurius6666/apiservice/src/pb/chatpb"
	"github.com/google/wire"
)

type Service interface {
	NewConversation(context.Context, string, string, []string, []string) (string, error)
	GetConversation(context.Context, string, []string) (*chatpb.ConversationPostResponse_Data, error)
	CloseConversation(context.Context, string) error
}

var Set = wire.NewSet(wire.Bind(new(Service), new(*ServiceGRPC)), wire.Struct(new(ServiceGRPC), "*"), ConnectClient)
