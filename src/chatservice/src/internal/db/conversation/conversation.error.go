package conversation

import "golang.org/x/xerrors"

var (
	prefix        = "conversation"
	ErrNotFound   = xerrors.Errorf("%s: record not found", prefix)
	ErrInsertFail = xerrors.Errorf("%s: insert failed", prefix)
)
