package user

import "golang.org/x/xerrors"

var (
	prefix        = "user"
	ErrNotFound   = xerrors.Errorf("%s: record not found", prefix)
	ErrInsertFail = xerrors.Errorf("%s: insert failed", prefix)
)
