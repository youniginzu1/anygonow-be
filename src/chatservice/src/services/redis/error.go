package redis

import "golang.org/x/xerrors"

var (
	prefix         = "redis"
	ErrKeyNotFound = xerrors.Errorf("%s: key not found", prefix)
	ErrInternal    = xerrors.Errorf("%s: internal error", prefix)
)
