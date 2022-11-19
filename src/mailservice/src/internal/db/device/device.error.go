package device

import "golang.org/x/xerrors"

var (
	prefix        = "device"
	ErrNotFound   = xerrors.Errorf("%s: record not found", prefix)
	ErrInsertFail = xerrors.Errorf("%s: insert failed", prefix)
)
