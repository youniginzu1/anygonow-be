package service

import "golang.org/x/xerrors"

var (
	prefix        = "service"
	ErrNotFound   = xerrors.Errorf("%s: record not found", prefix)
	ErrInsertFail = xerrors.Errorf("%s: insert failed", prefix)
)
