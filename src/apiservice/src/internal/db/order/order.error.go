package order

import "golang.org/x/xerrors"

var (
	prefix        = "order"
	ErrNotFound   = xerrors.Errorf("%s: record not found", prefix)
	ErrInsertFail = xerrors.Errorf("%s: insert failed", prefix)
)
