package cockroach

import (
	"context"

	"github.com/aqaurius6666/apiservice/src/internal/db/category"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

func applySearchCategory(db *gorm.DB, search *category.Search) *gorm.DB {
	if search.ID != uuid.Nil {
		db = db.Where(category.Category{
			BaseModel: database.BaseModel{
				ID: search.ID,
			},
		})
	}
	if search.CategoryIds != nil {
		db = db.Where(`"categories"."id"  = any (? :: uuid[])`, search.CategoryIds)
	}
	if search.Name != nil {
		db = db.Where(`UPPER("categories"."name") like UPPER(?)`, "%"+*search.Name+"%")
	}
	if search.Skip != 0 {
		db = db.Offset(search.Skip)
	}

	if search.Limit != 0 {
		db = db.Limit(search.Limit)
	}

	return db
}

func (u *ServerCDBRepo) SelectCategory(ctx context.Context, search *category.Search) (*category.Category, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.SelectCategory))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	r := category.Category{}
	if err := applySearchCategory(u.Db, search).WithContext(ctx).First(&r).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, xerrors.Errorf("%w", category.ErrNotFound)
		}
		return nil, err
	}
	return &r, nil
}

func (u *ServerCDBRepo) InsertCategory(ctx context.Context, value *category.Category) (*category.Category, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.InsertCategory))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := u.Db.WithContext(ctx).Create(value).Error; err != nil {
		err = xerrors.Errorf("%w", category.ErrInsertFail)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return value, nil
}

func (u *ServerCDBRepo) ListCategorys(ctx context.Context, search *category.Search) ([]*category.Category, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.ListCategorys))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	r := make([]*category.Category, 0)
	if err := applySearchCategory(u.Db, search).Joins(`left join "services" on "services"."category_id" = "categories"."id"`).
		Group(`"categories"."id", "categories"."name"`).
		Order(`COUNT("services"."category_id") DESC`).
		WithContext(ctx).Select(search.Fields).Find(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return r, nil
}

func (u *ServerCDBRepo) UpdateCategory(ctx context.Context, search *category.Search, value *category.Category) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.UpdateCategory))
	defer span.End()

	if err := applySearchCategory(u.Db, search).WithContext(ctx).Updates(value).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}

func (u *ServerCDBRepo) TotalCategory(ctx context.Context, search *category.Search) (*int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.TotalCategory))
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var r int64
	if err := applySearchCategory(u.Db, search).WithContext(ctx).Model(category.Category{}).
		Where(`"categories"."id" not in (?)`, subQueryForCheckCategoryNotInclude(u.Db, search)).Count(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &r, nil
}

func (u *ServerCDBRepo) DeleteCategory(ctx context.Context, value *category.Category) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.DeleteCategory))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := u.Db.WithContext(ctx).Delete(value).Error; err != nil {
		err = xerrors.Errorf("%w", category.ErrInsertFail)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}
func (u *ServerCDBRepo) ListCategoriesAdmin(ctx context.Context, search *category.Search) ([]*category.Category, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.ListCategoriesAdmin))
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	r := make([]*category.Category, 0)
	if err := applySearchCategory(u.Db, search).WithContext(ctx).
		Joins(`left join "services" on "categories"."id" = "services"."category_id" and "services"."status" = 0`).
		Joins(`left join groups "groups" on cast("categories"."id" as uuid) = any("groups"."category_ids" :: uuid[])`).
		Group(`"categories"."name", "categories"."id","groups"."fee","categories"."image_url"`).
		Where(`"categories"."id" not in (?)`, subQueryForCheckCategoryNotInclude(u.Db, search)).
		Select(search.Fields).
		Find(&r).Error; err != nil {

		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return r, nil
}

func subQueryForCheckCategoryNotInclude(db *gorm.DB, search *category.Search) *gorm.DB {
	db = db.Model(category.Category{})
	db = db.Select(`"categories"."id"`)
	if search.Query == c.QUERY_CATEGORY_ADMIN_ADVERTISE {
		db = db.Joins(`join "advertise_packages" on ("categories"."id" = any("advertise_packages"."categories" :: uuid[]))`)
	} else if search.Query == c.QUERY_CATEGORY_ADMIN_GROUP {
		db = db.Joins(`join "groups" on ("categories"."id" = any("groups"."category_ids" :: uuid[]))`)
	} else {
		db = db.Where(`false`)
	}
	return db
}
