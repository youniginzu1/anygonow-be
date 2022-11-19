package payment

import "golang.org/x/xerrors"

var (
	prefix        = "payment"
	ErrNotFound   = xerrors.Errorf("%s: record not found", prefix)
	ErrInsertFail = xerrors.Errorf("%s: insert failed", prefix)
)
