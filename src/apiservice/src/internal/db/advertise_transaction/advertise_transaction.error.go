package advertise_transaction

import "golang.org/x/xerrors"

var (
	prefix        = "advertise_transaction"
	ErrNotFound   = xerrors.Errorf("%s: record not found", prefix)
	ErrInsertFail = xerrors.Errorf("%s: insert failed", prefix)
)
