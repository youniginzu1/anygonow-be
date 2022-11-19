package cockroach

import (
	"context"

	"github.com/aqaurius6666/apiservice/src/internal/db/advertise_package"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

func applySearchAdvertisePackage(db *gorm.DB, search *advertise_package.Search) *gorm.DB {
	if search.ID != uuid.Nil {
		db = db.Where(advertise_package.AdvertisePackage{
			BaseModel: database.BaseModel{
				ID: search.ID,
			},
		})
	}

	return db
}

func (u *ServerCDBRepo) SelectAdvertisePackage(ctx context.Context, search *advertise_package.Search) (*advertise_package.AdvertisePackage, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.SelectAdvertisePackage))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r := advertise_package.AdvertisePackage{}
	if err := applySearchAdvertisePackage(u.Db, search).WithContext(ctx).First(&r).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, xerrors.Errorf("%w", advertise_package.ErrNotFound)
		}
		return nil, err
	}
	return &r, nil
}

func (u *ServerCDBRepo) InsertAdvertisePackage(ctx context.Context, value *advertise_package.AdvertisePackage) (*advertise_package.AdvertisePackage, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.InsertAdvertisePackage))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := u.Db.WithContext(ctx).Create(value).Error; err != nil {
		return nil, xerrors.Errorf("%w", advertise_package.ErrInsertFail)
	}
	return value, nil
}

func (u *ServerCDBRepo) UpdateAdvertisePackage(ctx context.Context, search *advertise_package.Search, value *advertise_package.AdvertisePackage) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.UpdateAdvertisePackage))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := applySearchAdvertisePackage(u.Db, search).WithContext(ctx).Updates(value).Error; err != nil {
		return xerrors.Errorf("%w", err)
	}
	return nil
}

func (u *ServerCDBRepo) TotalPackage(ctx context.Context, search *advertise_package.Search) (*int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.TotalPackage))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	var r int64
	db := u.Db
	if err := db.Table(`(?) as "t"`, GetSubAdvertisePackage(db, search)).
		Where(`UPPER(t.name) like UPPER(?)`, "%"+search.ServiceName+"%").WithContext(ctx).Count(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &r, nil
}

func (u *ServerCDBRepo) ListPackages(ctx context.Context, search *advertise_package.Search) ([]*advertise_package.AdvertisePackage, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.ListPackages))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r := make([]*advertise_package.AdvertisePackage, 0)

	db := u.Db
	if err := db.Table(`(?) as "t"`, GetSubAdvertisePackage(db, search)).
		Where(`UPPER(t.name) like UPPER(?)`, "%"+search.ServiceName+"%").WithContext(ctx).Find(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return r, nil
}

func GetSubAdvertisePackage(db *gorm.DB, search *advertise_package.Search) *gorm.DB {
	if search.Skip != 0 {
		db = db.Offset(search.Skip)
	}

	if search.Limit != 0 {
		db = db.Limit(search.Limit)
	}

	db = db.Table("advertise_packages")
	db = db.Select(search.Fields)
	db = db.Joins(`LEFT JOIN "categories" on "categories"."id" = any("advertise_packages"."categories" :: uuid[])`)
	db = db.Group(`"advertise_packages"."id", "advertise_packages"."name", "advertise_packages"."price", "advertise_packages"."banner_url", "advertise_packages"."description"`)
	return db
}

func (u *ServerCDBRepo) DeleteAdvertisePackage(ctx context.Context, value *advertise_package.AdvertisePackage) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.DeleteAdvertisePackage))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	if err := u.Db.WithContext(ctx).Delete(value).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}

func applyCheckExistedAdvertise(db *gorm.DB, search *advertise_package.Search) *gorm.DB {
	if search.ID != uuid.Nil {
		db = db.Where(`"advertise_packages"."id" != ?`, search.ID)
	}
	return db
}

func (u *ServerCDBRepo) CheckCateIdExistedAdvertise(ctx context.Context, search *advertise_package.Search, categoryId uuid.UUID) (*advertise_package.AdvertisePackage, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.CheckCateIdExistedAdvertise))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r := &advertise_package.AdvertisePackage{}

	if err := applyCheckExistedAdvertise(u.Db, search).WithContext(ctx).
		Where(`cast(? as uuid) = any("advertise_packages"."categories" :: uuid[])`, categoryId).
		Select(`"advertise_packages".*`).First(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return r, nil

}

func (u *ServerCDBRepo) ListAdvertiseDetails(ctx context.Context, search *advertise_package.Search) ([]*advertise_package.AdvertisePackage, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.ListAdvertiseDetails))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r := make([]*advertise_package.AdvertisePackage, 0)

	if err := applySearchAdvertisePackage(u.Db, search).WithContext(ctx).
		Joins(`left join "categories" on "categories"."id" = any ("advertise_packages"."categories" :: uuid[])`).
		Joins(`left join "businesses" on "businesses"."id" = (?)`, search.BusinessId).
		Group(`"advertise_packages"."id", "advertise_packages"."name", "advertise_packages"."price",
		"advertise_packages"."banner_url", "advertise_packages"."description","businesses"."zipcodes"`).
		Select(search.Fields).Find(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return r, nil

}

func (u *ServerCDBRepo) TotalAdvertiseDetail(ctx context.Context, search *advertise_package.Search) (*int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.TotalAdvertiseDetail))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var r int64

	if err := applySearchAdvertisePackage(u.Db, search).WithContext(ctx).Model(advertise_package.AdvertisePackage{}).Count(&r).
		Joins(`left join "categories" on "categories"."id" = any ("advertise_packages"."categories" :: uuid[])`).
		Joins(`left join "businesses" on "businesses"."id" = (?)`, search.BusinessId).
		Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &r, nil
}
