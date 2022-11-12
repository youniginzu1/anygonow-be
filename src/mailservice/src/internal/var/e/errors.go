package e

import (
	"golang.org/x/xerrors"
)

var (
	prefix                     = "mailservice"
	ErrAuthVerifySignatureFail = xerrors.Errorf("%s: verify signature fail", prefix)
	ErrInvalidActionType       = xerrors.Errorf("%s: invalid action type", prefix)
	ErrInvalidTimestamp        = xerrors.Errorf("%s: invalid timestamp", prefix)
	ErrIdInvalidFormat         = xerrors.Errorf("%s: invalid id", prefix)
	ErrOTPExpired              = xerrors.Errorf("%s: otp expired", prefix)
	ErrInvalidRequest          = xerrors.Errorf("%s: invalid request", prefix)
)
