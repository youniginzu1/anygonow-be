package otp

import "context"

type OtpRepo interface {
	SelectOtp(context.Context, *Search) (*Otp, error)
	InsertOtp(context.Context, *Otp) (*Otp, error)
	DeleteOTP(context.Context, *Search) error
	UpdateOTP(context.Context, *Search, *Otp) error
}
