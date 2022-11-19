package cockroach

import (
	"context"

	"github.com/aqaurius6666/apiservice/src/internal/db/payment"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

func applySearchPayment(db *gorm.DB, search *payment.Search) *gorm.DB {
	if search.ID != uuid.Nil {
		db = db.Where(payment.Payment{
			BaseModel: database.BaseModel{
				ID: search.ID,
			},
		})
	}
	if search.BusinessId != uuid.Nil {
		db = db.Where(payment.Payment{
			BusinessId: search.BusinessId,
		})
	}
	if search.CustomerId != nil {
		db = db.Where(payment.Payment{
			CustomerId: search.CustomerId,
		})
	}

	return db
}

func (u *ServerCDBRepo) SelectPayment(ctx context.Context, search *payment.Search) (*payment.Payment, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.SelectPayment))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r := payment.Payment{}
	if err := applySearchPayment(u.Db, search).WithContext(ctx).First(&r).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, xerrors.Errorf("%w", payment.ErrNotFound)
		}
		return nil, err
	}
	return &r, nil
}

func (u *ServerCDBRepo) InsertPayment(ctx context.Context, value *payment.Payment) (*payment.Payment, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.InsertPayment))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	if err := u.Db.WithContext(ctx).Create(value).Error; err != nil {
		err = xerrors.Errorf("%w", payment.ErrInsertFail)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return value, nil
}

func (u *ServerCDBRepo) ListPayments(ctx context.Context, search *payment.Search) ([]*payment.Payment, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.ListPayments))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r := make([]*payment.Payment, 0)
	if err := applySearchPayment(u.Db, search).WithContext(ctx).Find(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return r, nil
}

func (u *ServerCDBRepo) UpdatePayment(ctx context.Context, search *payment.Search, value *payment.Payment) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.UpdatePayment))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := applySearchPayment(u.Db, search).WithContext(ctx).Updates(value).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}
