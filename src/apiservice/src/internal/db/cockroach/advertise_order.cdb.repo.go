package cockroach

import (
	"context"
	"time"

	"github.com/aqaurius6666/apiservice/src/internal/db/advertise_order"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

func applySearchAdvertiseOrder(db *gorm.DB, search *advertise_order.Search) *gorm.DB {
	if search.ID != uuid.Nil {
		db = db.Where(advertise_order.AdvertiseOrder{
			BaseModel: database.BaseModel{
				ID: search.ID,
			},
		})
	}

	if search.AdvertiseOrder.BusinessId != uuid.Nil {
		db = db.Where(advertise_order.AdvertiseOrder{
			BusinessId: search.AdvertiseOrder.BusinessId,
		})
	}
	return db
}

func (u *ServerCDBRepo) SelectAdvertiseOrder(search *advertise_order.Search) (*advertise_order.AdvertiseOrder, error) {
	r := advertise_order.AdvertiseOrder{}
	if err := applySearchAdvertiseOrder(u.Db, search).First(&r).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, xerrors.Errorf("%w", advertise_order.ErrNotFound)
		}
		return nil, err
	}
	return &r, nil
}

func (u *ServerCDBRepo) InsertAdvertiseOrder(ctx context.Context, value *advertise_order.AdvertiseOrder) (*advertise_order.AdvertiseOrder, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.InsertAdvertiseOrder))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := u.Db.WithContext(ctx).Create(value).Error; err != nil {
		return nil, xerrors.Errorf("%w", advertise_order.ErrInsertFail)
	}
	return value, nil
}

func (u *ServerCDBRepo) UpdateAdvertiseOrder(search *advertise_order.Search, value *advertise_order.AdvertiseOrder) error {
	if err := applySearchAdvertiseOrder(u.Db, search).Updates(value).Error; err != nil {
		return xerrors.Errorf("%w", err)
	}
	return nil
}

func (u *ServerCDBRepo) TotalAdvertiseOrder(ctx context.Context, search *advertise_order.Search) (*int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.TotalAdvertiseOrder))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	var r int64
	if err := applySearchAdvertiseOrder(u.Db, search).WithContext(ctx).Model(advertise_order.AdvertiseOrder{}).Count(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &r, nil
}

func (u *ServerCDBRepo) ListAdvertiseOrders(ctx context.Context, search *advertise_order.Search) ([]*advertise_order.AdvertiseOrder, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.ListAdvertiseOrders))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	now := time.Now().UnixMilli()
	r := make([]*advertise_order.AdvertiseOrder, 0)
	if err := applySearchAdvertiseOrder(u.Db, search).Select(search.Fields).
		Joins(`left join "advertise_transactions" as a on "a"."id" = "advertise_orders"."advertise_transaction_id"`).
		Where(`"advertise_orders"."end_date" > ?`, now).
		WithContext(ctx).Find(&r).Error; err != nil {
		return nil, xerrors.Errorf("%w", err)
	}
	return r, nil
}

func (u *ServerCDBRepo) GetTotalOrderForBuyValidate(ctx context.Context, search *advertise_order.Search) (*int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.GetTotalOrderForBuyValidate))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	var r int64
	if err := applySearchAdvertiseOrder(u.Db, search).WithContext(ctx).Model(advertise_order.AdvertiseOrder{}).
		Joins(`join "services" on "advertise_orders"."service_id" = "services"."id"`).
		Joins(`join advertise_transactions on "advertise_orders"."advertise_transaction_id" = "advertise_transactions"."id"`).
		Where(`"services"."category_id" = ? AND "advertise_transactions"."zipcode" = ?`, search.CategoryId, search.Zipcode).
		Count(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &r, nil
}

func (u *ServerCDBRepo) TotalFeeAdvertise(ctx context.Context, search *advertise_order.Search) (*float64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.TotalFeeAdvertise))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	r := advertise_order.AdvertiseOrder{}
	if err := applySearchAdvertiseOrder(u.Db, search).WithContext(ctx).Model(advertise_order.AdvertiseOrder{}).
		Select(`SUM(price * ((end_date/1000 - start_date/1000) / 3600/ 24)) as "price"`).
		Joins(`join "advertise_transactions" on "advertise_orders"."advertise_transaction_id" = "advertise_transactions".id`).
		Group(`"advertise_orders"."business_id"`).Limit(1).
		Find(&r).
		Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return r.Price, nil
}
