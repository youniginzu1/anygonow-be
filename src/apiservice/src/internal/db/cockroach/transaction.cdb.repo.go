package cockroach

import (
	"context"

	"github.com/aqaurius6666/apiservice/src/internal/db/transaction"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

func applySearchTransaction(db *gorm.DB, search *transaction.Search) *gorm.DB {
	if search.BusinessId != uuid.Nil {
		db = db.Where(transaction.Transaction{
			BusinessId: search.BusinessId,
		})
	}
	if search.PaymentIntentId != nil {
		db = db.Where(transaction.Transaction{
			PaymentIntentId: search.PaymentIntentId,
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

	return db
}

func (u *ServerCDBRepo) ListTransactions(ctx context.Context, search *transaction.Search) ([]*transaction.Transaction, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.ListTransactions))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r := make([]*transaction.Transaction, 0)

	left, right := lib.GetTimeRange(*search.Query)

	if err := applySearchTransaction(u.Db, search).WithContext(ctx).
		Joins(`left join orders "orders" on "transactions"."order_id" = "orders"."id"`).
		Joins(`left join services "services" on "orders"."service_id" = "services"."id"`).
		Joins(`left join categories "categories" on "services"."category_id" = "categories"."id"`).
		Joins(`left join users "users" on "orders"."customer_id" = "users"."id"`).
		Joins(`left join businesses "businesses" on "orders"."business_id" = "businesses"."id"`).
		Where(`"transactions"."created_at" >= ? AND "transactions"."created_at" < ?`, left, right).
		Order(`"transactions"."created_at" DESC`).
		Select(search.Fields).Find(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return r, nil
}

func (u *ServerCDBRepo) TotalTransaction(ctx context.Context, search *transaction.Search) (*int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.TotalTransaction))
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var r int64
	left, right := lib.GetTimeRange(*search.Query)

	if err := applySearchTransaction(u.Db, search).
		Joins(`left join orders "orders" on "transactions"."order_id" = "orders"."id"`).
		Joins(`left join services "services" on "orders"."service_id" = "services"."id"`).
		Joins(`left join categories "categories" on "services"."category_id" = "categories"."id"`).
		Joins(`left join users "users" on "orders"."customer_id" = "users"."id"`).
		Joins(`left join businesses "businesses" on "orders"."business_id" = "businesses"."id"`).
		Where(`"transactions"."created_at" >= ? AND "transactions"."created_at" < ?`, left, right).
		WithContext(ctx).Model(transaction.Transaction{}).Count(&r).
		Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &r, nil
}

func (u *ServerCDBRepo) TotalFee(ctx context.Context, search *transaction.Search) (*int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.TotalFee))
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var r int64
	left, right := lib.GetTimeRange(*search.Query)
	if err := applySearchTransaction(u.Db, search).WithContext(ctx).Model(transaction.Transaction{}).Select(`coalesce(sum("transactions"."fee"),0) as "totalFee"`).
		Joins(`left join orders "orders" on "transactions"."order_id" = "orders"."id"`).
		Where(`"transactions"."created_at" >= ? AND "transactions"."created_at" < ?`, left, right).
		Scan(&r).
		Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &r, nil
}

func (u *ServerCDBRepo) InsertTransaction(ctx context.Context, value *transaction.Transaction) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.TotalFee))
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := u.Db.WithContext(ctx).Create(value).Error; err != nil {
		err = xerrors.Errorf("%w", transaction.ErrInsertFail)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}

func (u *ServerCDBRepo) UpdateTransaction(ctx context.Context, search *transaction.Search, value *transaction.Transaction) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.UpdateTransaction))
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := applySearchTransaction(u.Db, search).WithContext(ctx).Updates(value).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}
