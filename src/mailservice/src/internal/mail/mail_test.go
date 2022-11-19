package mail

import (
	"bytes"
	"os"
	"testing"
	"text/template"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	AWS_USERNAME  = os.Getenv("CONFIG_MAIL_USERNAME")
	AWS_PASS      = os.Getenv("CONFIG_MAIL_PASSWORD")
	AWS_SMTP_HOST = os.Getenv("CONFIG_MAIL_SMTP_HOST")
	AWS_SMTP_PORT = os.Getenv("CONFIG_MAIL_SMTP_PORT")
	SMTP_SENDER   = os.Getenv("CONFIG_MAIL_SMTP_SENDER")
)

func TestSendMail(t *testing.T) {
	tmpl, err := template.New("test").ParseFiles("mail.tpl")
	assert.Nil(t, err)
	bz := new(bytes.Buffer)
	to := "anygonow@gmail.com"
	err = tmpl.ExecuteTemplate(bz, "mail.tpl", struct {
		UserEmail    string
		EmailSender  string
		Subject      string
		Boundary     string
		Text1        string
		Text2        string
		CompanyEmail string
		Link         string
		ButtonText   string
	}{
		UserEmail:    to,
		EmailSender:  SMTP_SENDER,
		Subject:      "This is test subject",
		Boundary:     uuid.NewString(),
		Text1:        "This is test text 1",
		Text2:        "This is test text 2",
		CompanyEmail: "support@anygonow.com",
		Link:         "https://uet.vnu.edu.vn",
		ButtonText:   "Button",
	})
	assert.Nil(t, err, err)
	client := NewServiceV2(MailUsername(AWS_USERNAME), MailPassword(AWS_PASS), SMTPHost(AWS_SMTP_HOST), SMTPPort(AWS_SMTP_PORT), SMTPSender(SMTP_SENDER))
	err = client.SendMail(to, bz.Bytes())
	assert.Nil(t, err)
	if err != nil {
		t.Log(err)
	}
}
