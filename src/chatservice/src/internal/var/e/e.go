package e

import (
	"golang.org/x/xerrors"
)

var (
	prefix                     = "chatservice"
	ErrMissingField            = func(a string) error { return xerrors.Errorf("%s: missing field: %s", prefix, a) }
	ErrUserInactive            = xerrors.Errorf("%s: user inactive", prefix)
	ErrIdInvalidFormat         = xerrors.Errorf("%s: id invalid format", prefix)
	ErrBodyInvalid             = xerrors.Errorf("%s: body invalid", prefix)
	ErrMissingBody             = xerrors.Errorf("%s: missing body", prefix)
	ErrMissingCertificate      = xerrors.Errorf("%s: missing certificate", prefix)
	ErrExpectedBody            = xerrors.Errorf("%s: expected body", prefix)
	ErrNoPhoneAvailable        = xerrors.Errorf("%s: no phone available", prefix)
	ErrNoConversationAvailable = xerrors.Errorf("%s: unable to release conversation phone pool", prefix)
	ErrOpenConversation        = xerrors.Errorf("%s: open conversation", prefix)
	ErrBuyPhoneFail            = xerrors.Errorf("%s: unable to buy new phone number", prefix)
	ErrMaxPhoneNumbers         = xerrors.Errorf("%s: unable to buy more phone number", prefix)
)
