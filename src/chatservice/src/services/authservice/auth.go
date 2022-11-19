package authservice

import (
	"context"

	"github.com/aqaurius6666/chatservice/src/internal/var/c"
	"github.com/google/wire"
)

var Set = wire.NewSet(wire.Bind(new(Service), new(ServiceGRPC)), wire.Struct(new(ServiceGRPC), "*"), ConnectClient)

type Service interface {
	CheckAuth(ctx context.Context, header, body []byte, method string) (string, c.ROLE, error)
	// GetCredential(ctx context.Context, data string) (id, pub, priv string, err error)
	// Register(ctx context.Context, mail, phone, pub, enpriv string, role c.ROLE) (otpId string, err error)
	// RegisterOTP(ctx context.Context, otpId string) (id, mail, phone string, role c.ROLE, err error)
	// ForgotPasswordOTP(ctx context.Context, otpId, pub, enc string) (err error)
	// ForgotPassword(ctx context.Context, mail string) (otpId string, err error)
	// VerifyOTP(ctx context.Context, otpId, otp string) (c.OTP_TYPE, error)
	// ChangePassword(ctx context.Context, id, pub, enc string) error
	// BanUser(ctx context.Context, id string, status bool) (string, c.ROLE, error)
	// DeleteUser(ctx context.Context, id string) (string, c.ROLE, error)
	// CheckCredential(ctx context.Context, identifier string) (bool, error)
	// ChangeMail(ctx context.Context, id string, mail string) (string, error)
	// ChangeMailOTP(ctx context.Context, otpId string) (string, string, c.ROLE, error)
}
