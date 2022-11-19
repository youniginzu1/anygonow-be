package business

import "golang.org/x/xerrors"

var (
	prefix        = "business"
	ErrNotFound   = xerrors.Errorf("%s: record not found", prefix)
	ErrInsertFail = xerrors.Errorf("%s: insert failed", prefix)
)
