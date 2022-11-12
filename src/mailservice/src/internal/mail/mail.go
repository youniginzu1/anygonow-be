package mail

import (
	"github.com/google/wire"
)

var (
	ServiceSet = wire.NewSet(wire.Bind(new(Service), new(*ServiceV2)), NewServiceV2)
)

type (
	MailUsername string
	MailPassword string
	SMTPHost     string
	SMTPPort     string
	SMTPSender   string
)
type Service interface {
	SendMails(to []string, message []byte) error
	SendMail(to string, message []byte) error
}
