//go:build wireinject
// +build wireinject

package db

import (
	"context"

	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

func NewMockRepo(ctx context.Context, logger *logrus.Logger, dsn DBDsn) (ServerRepo, error) {
	wire.Build(
		ServerRepoSet,
	)

	return nil, nil
}
