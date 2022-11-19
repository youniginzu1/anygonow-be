package transaction

import "golang.org/x/xerrors"

var (
	prefix        = "transaction"
	ErrNotFound   = xerrors.Errorf("%s: record not found", prefix)
	ErrInsertFail = xerrors.Errorf("%s: insert failed", prefix)
)
