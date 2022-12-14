package feedback

import "golang.org/x/xerrors"

var (
	prefix        = "feedback"
	ErrNotFound   = xerrors.Errorf("%s: record not found", prefix)
	ErrInsertFail = xerrors.Errorf("%s: insert failed", prefix)
)
