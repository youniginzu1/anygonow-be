package cockroach

import (
	"context"

	"github.com/aqaurius6666/apiservice/src/internal/db/advertise_transaction"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

func applySearchAdvertiseTransaction(db *gorm.DB, search *advertise_transaction.Search) *gorm.DB {
	if search.ID != uuid.Nil {
		db = db.Where(advertise_transaction.AdvertiseTransaction{
			BaseModel: database.BaseModel{
				ID: search.ID,
			},
		})
	}
	return db
}

func (u *ServerCDBRepo) SelectAdvertiseTransaction(search *advertise_transaction.Search) (*advertise_transaction.AdvertiseTransaction, error) {
	r := advertise_transaction.AdvertiseTransaction{}
	if err := applySearchAdvertiseTransaction(u.Db, search).First(&r).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, xerrors.Errorf("%w", advertise_transaction.ErrNotFound)
		}
		return nil, err
	}
	return &r, nil
}

func (u *ServerCDBRepo) InsertAdvertiseTransaction(ctx context.Context, value *advertise_transaction.AdvertiseTransaction) (*advertise_transaction.AdvertiseTransaction, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.InsertAdvertiseTransaction))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)

	defer cancel()
	if err := u.Db.WithContext(ctx).Create(value).Error; err != nil {
		return nil, xerrors.Errorf("%w", advertise_transaction.ErrInsertFail)
	}
	return value, nil
}

func (u *ServerCDBRepo) UpdateAdvertiseTransaction(search *advertise_transaction.Search, value *advertise_transaction.AdvertiseTransaction) error {
	if err := applySearchAdvertiseTransaction(u.Db, search).Updates(value).Error; err != nil {
		return xerrors.Errorf("%w", err)
	}
	return nil
}
