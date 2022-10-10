package role

import "golang.org/x/xerrors"

var (
	prefix        = "role"
	ErrNotFound   = xerrors.Errorf("%s: record not found", prefix)
	ErrInsertFail = xerrors.Errorf("%s: insert failed", prefix)
)
