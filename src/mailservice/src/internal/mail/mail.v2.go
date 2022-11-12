package mail

import (
	"fmt"
	"net/smtp"
)

type ServiceV2 struct {
	auth       smtp.Auth
	from       string
	smtpServer string
}

func NewServiceV2(user MailUsername, pass MailPassword, host SMTPHost, port SMTPPort, mail SMTPSender) *ServiceV2 {
	smtpServer := fmt.Sprintf("%s:%s", host, port)
	auth := smtp.PlainAuth("", string(user), string(pass), string(host))
	return &ServiceV2{
		auth:       auth,
		from:       string(mail),
		smtpServer: smtpServer,
	}
}

func (s *ServiceV2) SendMails(to []string, message []byte) error {
	return smtp.SendMail(s.smtpServer, s.auth, s.from, to, message)
}

func (s *ServiceV2) SendMail(to string, message []byte) error {
	return smtp.SendMail(s.smtpServer, s.auth, s.from, []string{to}, message)
}
