package cockroach

import (
	"context"

	"github.com/aqaurius6666/apiservice/src/internal/db/service"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

func applySearchService(db *gorm.DB, search *service.Search) *gorm.DB {
	if search.ID != uuid.Nil {
		db = db.Where(service.Service{
			BaseModel: database.BaseModel{
				ID: search.ID,
			},
		})
	}
	if search.BusinessId != uuid.Nil {
		db = db.Where(service.Service{
			BusinessId: search.BusinessId,
		})
	}
	if search.CategoryId != uuid.Nil {
		db = db.Where(service.Service{
			CategoryId: search.CategoryId,
		})
	}
	if search.Status != nil {
		db = db.Where(service.Service{
			Status: search.Status,
		})
	}
	if search.CategoryIds != nil {
		db = db.Where(`"services"."category_id" = any (? :: uuid[])`, search.CategoryIds)
	}
	return db
}

func (u *ServerCDBRepo) SelectService(ctx context.Context, search *service.Search) (*service.Service, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.SelectService))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	r := service.Service{}
	if err := applySearchService(u.Db, search).WithContext(ctx).First(&r).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = service.ErrNotFound
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &r, nil
}

func (u *ServerCDBRepo) InsertService(ctx context.Context, value *service.Service) (*service.Service, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.InsertService))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := u.Db.WithContext(ctx).Create(value).Error; err != nil {
		err = xerrors.Errorf("%w", service.ErrInsertFail)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return value, nil
}

func (u *ServerCDBRepo) ListServices(ctx context.Context, search *service.Search) ([]*service.Service, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.ListServices))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r := make([]*service.Service, 0)
	if err := applySearchService(u.Db, search).
		WithContext(ctx).
		Joins(`left join "categories" on "categories"."id" = "services"."category_id"`).
		Joins(`left join "orders" on "orders"."service_id" = "services"."id" and "orders"."status" = ?`, c.ORDER_STATUS_COMPLETED).
		Group(`"services"."id", "categories"."name", "categories"."image_url", "categories"."id"`).
		Select(search.Fields).Find(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return r, nil
}

func (u *ServerCDBRepo) UpdateService(ctx context.Context, search *service.Search, value *service.Service) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.UpdateService))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	if err := applySearchService(u.Db, search).WithContext(ctx).Updates(value).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}

func (u *ServerCDBRepo) FirstOrInsertService(ctx context.Context, search *service.Search, value *service.Service) (*service.Service, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.FirstOrInsertService))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	if err := applySearchService(u.Db, search).WithContext(ctx).FirstOrCreate(value).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return value, nil
}
