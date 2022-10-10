package user

import "golang.org/x/xerrors"

var (
	prefix                  = "user"
	ErrNotFound             = xerrors.Errorf("%s: record not found", prefix)
	ErrInsertFail           = xerrors.Errorf("%s: insert failed", prefix)
	ErrUsernameExisted      = xerrors.Errorf("%s: username existed", prefix)
	ErrEmailExisted         = xerrors.Errorf("%s: email existed", prefix)
	ErrPhoneExisted         = xerrors.Errorf("%s: phone existed", prefix)
	ErrReferralCodeNotFound = xerrors.Errorf("referral code not found")
)
