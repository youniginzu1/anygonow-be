package mailservice

import (
	"context"

	"github.com/google/wire"
)

type Service interface {
	SubscribeNotification(ctx context.Context, userId string, deviceId string) error
	UnsubscribeNotification(ctx context.Context, userId string, deviceId string) error
	SendNotification(ctx context.Context, to string, title, body string, message string) error
}

var Set = wire.NewSet(wire.Bind(new(Service), new(*ServiceGRPC)), wire.Struct(new(ServiceGRPC), "*"), ConnectClient)
