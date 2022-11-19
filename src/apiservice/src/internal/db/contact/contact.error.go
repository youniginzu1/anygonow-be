package contact

import "golang.org/x/xerrors"

var (
	prefix        = "contact"
	ErrNotFound   = xerrors.Errorf("%s: record not found", prefix)
	ErrInsertFail = xerrors.Errorf("%s: insert failed", prefix)
)
