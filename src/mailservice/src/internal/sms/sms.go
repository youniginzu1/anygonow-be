package sms

import "context"

type Service interface {
	SendSMS(ctx context.Context, to string, msg string) error
}

func NewService() (Service, error) {
	return NewSNSService()
}
