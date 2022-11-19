package phone_pool

import "golang.org/x/xerrors"

var (
	prefix               = "phone_pool"
	ErrNotFound          = xerrors.Errorf("%s: record not found", prefix)
	ErrInsertFail        = xerrors.Errorf("%s: insert failed", prefix)
	ErrFirstOrCreateFail = xerrors.Errorf("%s: first or create failed", prefix)
)
