package group

import "golang.org/x/xerrors"

var (
	prefix        = "category"
	ErrNotFound   = xerrors.Errorf("%s: record not found", prefix)
	ErrInsertFail = xerrors.Errorf("%s: insert failed", prefix)
)
