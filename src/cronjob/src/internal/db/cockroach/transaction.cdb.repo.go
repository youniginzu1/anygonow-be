package cockroach

import (
	"context"

	"github.com/aqaurius6666/cronjob/src/internal/db/transaction"
	"github.com/aqaurius6666/cronjob/src/internal/lib"
	"github.com/aqaurius6666/cronjob/src/internal/var/c"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

func applySearchTransaction(db *gorm.DB, search *transaction.Search) *gorm.DB {
	if search.Now != nil {
		db = db.Where(`"transactions"."created_at" < ?`, *search.Now)
	}
	if search.BusinessId != uuid.Nil {
		db = db.Where(transaction.Transaction{
			BusinessId: search.BusinessId,
		})
	}
	if search.IsPaid != nil {
		db = db.Where(transaction.Transaction{
			IsPaid: search.IsPaid,
		})
	}
	if search.IsFree != nil {
		db = db.Where(transaction.Transaction{
			IsFree: search.IsFree,
		})
	}
	if search.BIds != nil {
		db = db.Where(`"transactions"."business_id" = any (? :: uuid[])`, search.BIds)
	}
	return db
}

func (u *ServerCDBRepo) GetTotalMoneyOfBusinessesPreviousMonth(ctx context.Context, search *transaction.Search) ([]*transaction.Transaction, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.GetTotalMoneyOfBusinessesPreviousMonth))
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r := make([]*transaction.Transaction, 0)
	if err := applySearchTransaction(u.Db, search).WithContext(ctx).Select(search.Fields).
		Joins(`join "payments" on "transactions"."business_id" = "payments"."business_id"`).
		Group(`"payments"."customer_id", "transactions"."business_id", "payments"."payment_method_id"`).Find(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}

	return r, nil
}

func (u *ServerCDBRepo) UpdateTransaction(ctx context.Context, search *transaction.Search, value *transaction.Transaction) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.UpdateTransaction))
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := applySearchTransaction(u.Db, search).WithContext(ctx).Updates(value).Error; err != nil {
		err = xerrors.Errorf("%w", transaction.ErrInsertFail)
		lib.RecordError(span, err)
		return err
	}
	return nil
}
