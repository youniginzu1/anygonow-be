package advertise_package

import "golang.org/x/xerrors"

var (
	prefix        = "advertise_package"
	ErrNotFound   = xerrors.Errorf("%s: record not found", prefix)
	ErrInsertFail = xerrors.Errorf("%s: insert failed", prefix)
)
