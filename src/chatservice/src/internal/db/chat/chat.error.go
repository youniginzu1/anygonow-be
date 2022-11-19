package chat

import "golang.org/x/xerrors"

var (
	prefix        = "chat"
	ErrNotFound   = xerrors.Errorf("%s: record not found", prefix)
	ErrInsertFail = xerrors.Errorf("%s: insert failed", prefix)
)
