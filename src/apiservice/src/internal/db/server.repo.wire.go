package db

import (
	"net/url"

	"github.com/aqaurius6666/apiservice/src/internal/db/cockroach"
	"github.com/google/wire"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"golang.org/x/xerrors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var ServerRepoSet = wire.NewSet(ConnectGorm, wire.Bind(new(ServerRepo), new(*cockroach.ServerCDBRepo)), cockroach.CDBRepoSet)

var gormConfig = &gorm.Config{
	SkipDefaultTransaction:                   true,
	DisableAutomaticPing:                     false,
	PrepareStmt:                              true,
	DisableForeignKeyConstraintWhenMigrating: true,
}

func ConnectGorm(dsn DBDsn) (db *gorm.DB, err error) {
	uri, err := url.Parse(string(dsn))
	if err != nil {
		return nil, xerrors.Errorf("could not parse DB URI: %w", err)
	}
	switch uri.Scheme {
	case "in-memory":
		return nil, xerrors.Errorf("Not implemented!")
	case "postgresql":
		db, err = gorm.Open(postgres.Open(string(dsn)), gormConfig)
		if err != nil {
			return nil, xerrors.Errorf("%w", err)
		}
	case "postgres":
		db, err = gorm.Open(postgres.Open(string(dsn)), gormConfig)
		if err != nil {
			return nil, xerrors.Errorf("%w", err)
		}
	default:
		return nil, xerrors.Errorf("unsupported DB URI scheme: %q", uri.Scheme)
	}

	db.Use(otelgorm.NewPlugin())
	return
}
