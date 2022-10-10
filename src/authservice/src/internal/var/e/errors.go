package e

import (
	"golang.org/x/xerrors"
)

var (
	prefix                     = "authservice"
	ErrAuthParseModelFail      = xerrors.Errorf("%s: parse model fail", prefix)
	ErrAuthExpired             = xerrors.Errorf("%s: certificate expired", prefix)
	ErrAuthVerifySignatureFail = xerrors.Errorf("%s: verify signature fail", prefix)
	ErrInvalidActionType       = xerrors.Errorf("%s: invalid action type", prefix)
	ErrInvalidTimestamp        = xerrors.Errorf("%s: invalid timestamp", prefix)
	ErrIdInvalidFormat         = xerrors.Errorf("%s: invalid id", prefix)
	ErrOTPExpired              = xerrors.Errorf("%s: otp expired", prefix)
	ErrInvalidError            = xerrors.Errorf("%s: invalid error", prefix)
	ErrOTPNotMatch             = xerrors.Errorf("%s: otp not match", prefix)
	ErrOTPInvalid              = xerrors.Errorf("%s: otp invalid", prefix)
	ErrDeleteActiveUser        = xerrors.Errorf("%s: delete active user", prefix)
	ErrUserInactive            = xerrors.Errorf("%s: user inactive", prefix)
	ErrOTPSpam                 = xerrors.Errorf("%s: spam", prefix)
)
