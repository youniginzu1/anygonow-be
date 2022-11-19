package advertise_order

import "golang.org/x/xerrors"

var (
	prefix        = "advertise_order"
	ErrNotFound   = xerrors.Errorf("%s: record not found", prefix)
	ErrInsertFail = xerrors.Errorf("%s: insert failed", prefix)
)
