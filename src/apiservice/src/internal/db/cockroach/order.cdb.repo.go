package cockroach

import (
	"context"
	"time"

	"github.com/aqaurius6666/apiservice/src/internal/db/order"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/aqaurius6666/go-utils/utils"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

func applySearchOrder(db *gorm.DB, search *order.Search) *gorm.DB {
	if search.ID != uuid.Nil {
		db = db.Where(order.Order{
			BaseModel: database.BaseModel{
				ID: search.ID,
			},
		})
	}
	if search.CustomerId != uuid.Nil {
		db = db.Where(order.Order{
			CustomerId: search.CustomerId,
		})
	}
	if search.CustomerZipcode != nil {
		db = db.Where(order.Order{
			CustomerZipcode: search.CustomerZipcode,
		})
	}
	if search.Status != nil {
		db = db.Where(order.Order{
			Status: search.Status,
		})
	}
	if search.UserId != uuid.Nil {
		db = db.Where(`"orders"."customer_id" = ? or "orders"."business_id" = ?`, search.UserId, search.UserId)
	}
	if search.CategoryId != uuid.Nil {
		db = db.Where(`"categories"."id" = ?`, search.CategoryId)
	}
	if search.ServiceId != uuid.Nil {
		db = db.Where(order.Order{
			ServiceId: search.ServiceId,
		})
	}
	if search.Skip != 0 {
		db = db.Offset(search.Skip)
	}

	if search.Limit != 0 {
		db = db.Limit(search.Limit)
	}

	if search.BusinessId != uuid.Nil {
		db = db.Where(order.Order{
			BusinessId: search.BusinessId,
		})
	}

	if search.OIds != nil {
		db = db.Where(`"orders"."id"  = any (? :: uuid[])`, search.OIds)
	}

	return db
}

func (u *ServerCDBRepo) SelectOrder(ctx context.Context, search *order.Search) (*order.Order, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.SelectOrder))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r := order.Order{}
	if err := applySearchOrder(u.Db, search).WithContext(ctx).First(&r).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, xerrors.Errorf("%w", order.ErrNotFound)
		}
		return nil, err
	}
	return &r, nil
}

func (u *ServerCDBRepo) InsertOrder(ctx context.Context, value *order.Order) (*order.Order, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.InsertOrder))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	if err := u.Db.WithContext(ctx).Create(value).Error; err != nil {
		err = xerrors.Errorf("%w", order.ErrInsertFail)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return value, nil
}

func (u *ServerCDBRepo) ListOrders(ctx context.Context, search *order.Search) ([]*order.Order, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.ListOrders))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r := make([]*order.Order, 0)
	if err := applySearchOrder(u.Db, search).WithContext(ctx).
		Joins(`left join services "services" on "orders"."service_id" = "services"."id"`).
		Joins(`left join categories "categories" on "services"."category_id" = "categories"."id"`).
		Joins(`left join users "users" on "orders"."customer_id" = "users"."id"`).
		Joins(`left join businesses "businesses" on "orders"."business_id" = "businesses"."id"`).
		Joins(`left join groups "groups" on cast("categories"."id" as uuid) = any("groups"."category_ids" :: uuid[])`).
		Order("orders.updated_at DESC").
		Select(search.Fields).Find(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return r, nil
}

func (u *ServerCDBRepo) UpdateOrder(ctx context.Context, search *order.Search, value *order.Order) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.UpdateOrder))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := applySearchOrder(u.Db, search).WithContext(ctx).Updates(value).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}

func (u *ServerCDBRepo) UpdateOrderStatusIfExpireTime(ctx context.Context) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.UpdateOrderStatusIfExpireTime))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var now = time.Now().UnixMilli()
	if err := u.Db.WithContext(ctx).Table("orders").
		Where(`"orders"."status" = ? AND "orders"."end_date" < ?`, c.ORDER_STATUS_PENDING, now).
		Updates(&order.Order{
			Status: utils.Int32Ptr(int32(c.ORDER_STATUS_REJECTED)),
		}).
		Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}

func (u *ServerCDBRepo) UpdateOrdersCancelledStatusByUser(ctx context.Context, search *order.Search, value *order.Order) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.UpdateOrdersCancelledStatusByUser))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := applySearchOrder(u.Db, search).WithContext(ctx).
		Where(`"orders"."status" = ? or "orders"."status" = ?`, c.ORDER_STATUS_PENDING, c.ORDER_STATUS_CONNECTED).
		Updates(value).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}

func (u *ServerCDBRepo) TotalOrder(ctx context.Context, search *order.Search) (*int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.TotalOrder))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	var r int64
	if err := applySearchOrder(u.Db, search).Model(order.Order{}).
		Joins(`left join services "services" on "orders"."service_id" = "services"."id"`).
		Joins(`left join categories "categories" on "services"."category_id" = "categories"."id"`).
		WithContext(ctx).Count(&r).Error; err != nil {
		return nil, xerrors.Errorf("%w", err)
	}
	return &r, nil
}

func (u *ServerCDBRepo) ListProjects(ctx context.Context, search *order.Search) ([]*order.Order, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.ListProjects))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r := make([]*order.Order, 0)
	if err := SubQueryProjects(u.Db, search).WithContext(ctx).Find(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return r, nil
}

func (u *ServerCDBRepo) TotalProjects(ctx context.Context, search *order.Search) (*int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.TotalProjects))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	var r int64
	if err := SubQueryProjects(u.Db, search).WithContext(ctx).Model(order.Order{}).Count(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &r, nil
}

func SubQueryProjects(db *gorm.DB, search *order.Search) *gorm.DB {
	if search.CustomerId != uuid.Nil {
		db = db.Where(order.Order{
			CustomerId: search.CustomerId,
		})
	}
	if search.Skip != 0 {
		db = db.Offset(search.Skip)
	}

	if search.Limit != 0 {
		db = db.Limit(search.Limit)
	}
	db = db.
		Joins(`left join services "services" on "orders"."service_id" = "services"."id"`).
		Joins(`left join categories on "categories"."id" = "services"."category_id"`).
		Group(`"categories"."id", "orders"."customer_zipcode"`).
		Order(`"categories"."name" asc`).
		Where(`"orders"."status" = ? or "orders"."status" = ?`, c.ORDER_STATUS_PENDING, c.ORDER_STATUS_CONNECTED).
		Select(search.Fields)
	return db
}

func (u *ServerCDBRepo) CancelProject(ctx context.Context, search *order.Search, value *order.Order) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.CancelProject))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	subQuery := getSubQueryForCancelProject(u.Db, search)
	if err := u.Db.WithContext(ctx).Where(`"orders"."id" IN (?)`, subQuery).
		Updates(value).
		Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}

func (u *ServerCDBRepo) ListBusinessesAlreadyOrdered(ctx context.Context, search *order.Search) ([]*order.Order, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.ListBusinessesAlreadyOrdered))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r := make([]*order.Order, 0)
	if err := applySearchOrder(u.Db, search).
		WithContext(ctx).
		Where(`"orders"."status" = ? or "orders"."status" = ?`, c.ORDER_STATUS_PENDING, c.ORDER_STATUS_CONNECTED).
		Joins(`left join "services" "services" on "orders"."service_id" = "services"."id"`).
		Joins(`left join "categories" "categories" on "services"."category_id" = "categories"."id"`).
		Select(search.Fields).
		Find(&r).
		Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, xerrors.Errorf("%w", order.ErrNotFound)
		}
		return nil, err
	}
	return r, nil
}

func getSubQueryForCancelProject(db *gorm.DB, search *order.Search) *gorm.DB {
	db = db.Table("orders")
	db = db.Joins(`JOIN "services" on "orders"."service_id" = "services"."id"`)
	db = db.Joins(`JOIN "categories" on "services"."category_id" = "categories"."id" and "categories"."id" = ?`, search.ServiceId)
	db = db.Where(`"orders"."customer_id" = ? AND "orders"."customer_zipcode" = ?`, search.CustomerId, search.CustomerZipcode)
	db = db.Select(`"orders"."id"`)
	return db
}
