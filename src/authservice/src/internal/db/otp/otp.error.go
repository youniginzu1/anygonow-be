package otp

import "golang.org/x/xerrors"

var (
	prefix        = "otp"
	ErrNotFound   = xerrors.Errorf("%s: record not found", prefix)
	ErrInsertFail = xerrors.Errorf("%s: insert failed", prefix)
	ErrDeleteFail = xerrors.Errorf("%s: delete failed", prefix)
)
