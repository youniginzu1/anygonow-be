package mailservice

import (
	"context"

	"github.com/google/wire"
)

var MailServiceSet = wire.NewSet(wire.Bind(new(Service), new(ServiceGRPC)), wire.Struct(new(ServiceGRPC), "*"), ConnectClient)

type Service interface {
	SendMail(ctx context.Context, to string, msg []byte) error
	SendMails(ctx context.Context, to []string, msg []byte) error
}
