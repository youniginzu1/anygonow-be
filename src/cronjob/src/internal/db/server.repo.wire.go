package db

import (
	"net/url"

	"github.com/aqaurius6666/cronjob/src/internal/db/cockroach"
	"github.com/google/wire"
	"golang.org/x/xerrors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var ServerRepoSet = wire.NewSet(ConnectGorm, wire.Bind(new(ServerRepo), new(*cockroach.ServerCDBRepo)), cockroach.CDBRepoSet)
var gormConfig = &gorm.Config{
	SkipDefaultTransaction:                   true,
	DisableAutomaticPing:                     true,
	DisableForeignKeyConstraintWhenMigrating: true,
}

func ConnectGorm(dsn DBDsn) (*gorm.DB, error) {
	uri, err := url.Parse(string(dsn))
	if err != nil {
		return nil, xerrors.Errorf("could not parse DB URI: %w", err)
	}
	switch uri.Scheme {
	case "in-memory":
		return nil, xerrors.Errorf("Not implemented!")
	case "postgresql":
		return gorm.Open(postgres.Open(string(dsn)), gormConfig)
	case "postgres":
		return gorm.Open(postgres.Open(string(dsn)), gormConfig)
	default:
		return nil, xerrors.Errorf("unsupported DB URI scheme: %q", uri.Scheme)
	}
}
