package notification

import "golang.org/x/xerrors"

var (
	prefix        = "notification"
	ErrNotFound   = xerrors.Errorf("%s: record not found", prefix)
	ErrInsertFail = xerrors.Errorf("%s: insert failed", prefix)
	ErrFirstOrCreateFail = xerrors.Errorf("%s: first or create failed", prefix)
)
