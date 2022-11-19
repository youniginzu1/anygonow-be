package sms

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func MockSNS() (Service, error) {

	os.Setenv("AWS_ACCESS_KEY_ID", "AKIAV4NSU3G25SRXVXU5")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "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")
	os.Setenv("AWS_REGION", "us-east-1")
	client, err := NewService()
	return client, err
}

func TestSendSMS(t *testing.T) {

	client, err := MockSNS()
	assert.Nil(t, err)
	err = client.SendSMS(context.Background(), "+84384511909", "hello from code")
	assert.Nil(t, err)
}
