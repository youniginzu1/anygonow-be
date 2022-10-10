package mailservice

import (
	"context"
	"testing"

	"github.com/aqaurius6666/authservice/src/internal/lib/template"
	"github.com/stretchr/testify/assert"
)

func TestSendMail(t *testing.T) {
	c, err := NewMailService()
	assert.Nil(t, err)
	err = c.SendMail(context.Background(), "aqaurius6666+1@gmail.com", []byte("hello"))
	assert.Nil(t, err)
}

func TestSendMailRegister(t *testing.T) {
	c, err := NewMailService()
	assert.Nil(t, err, err)
	toMail := "aqaurius1910@gmail.com"
	otp := "12123"
	otpId := "2312313123"
	tmp := template.RegisterMailTemplate(toMail, otp, otpId)
	assert.NotEqual(t, []byte("error mail"), tmp, "error mail")
	err = c.SendMail(context.Background(), toMail, tmp)
	assert.Nil(t, err, err)
}
