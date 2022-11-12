package mail

import (
	"fmt"
	"net/smtp"
)

type ServiceV1 struct {
	auth       smtp.Auth
	from       string
	smtpServer string
}

func NewServiceV1(user MailUsername, pass MailPassword, host SMTPHost, port SMTPPort) *ServiceV1 {
	smtpServer := fmt.Sprintf("%s:%s", host, port)
	auth := smtp.PlainAuth("", string(user), string(pass), string(host))
	return &ServiceV1{
		auth:       auth,
		from:       string(user),
		smtpServer: smtpServer,
	}
}

func (s *ServiceV1) SendMails(to []string, message []byte) error {
	return smtp.SendMail(s.smtpServer, s.auth, s.from, to, message)
}

func (s *ServiceV1) SendMail(to string, message []byte) error {
	fmt.Printf("message: %v\n", string(message))
	return smtp.SendMail(s.smtpServer, s.auth, s.from, []string{to}, message)
}
