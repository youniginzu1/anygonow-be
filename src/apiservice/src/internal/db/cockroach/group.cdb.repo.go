package cockroach

import (
	"context"

	"github.com/aqaurius6666/apiservice/src/internal/db/category"
	"github.com/aqaurius6666/apiservice/src/internal/db/group"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

func applySearchGroup(db *gorm.DB, search *group.Search) *gorm.DB {
	if search.ID != uuid.Nil {
		db = db.Where(group.Group{
			BaseModel: database.BaseModel{
				ID: search.ID,
			},
		})
	}
	if search.Name != nil {
		db = db.Where(group.Group{
			Name: search.Name,
		})
	}
	if search.Skip != 0 {
		db = db.Offset(search.Skip)
	}

	if search.Limit != 0 {
		db = db.Limit(search.Limit)
	}
	return db
}

func applySelectGroup(db *gorm.DB, search *group.Search) *gorm.DB {
	if search.ID != uuid.Nil {
		db = db.Where(group.Group{
			BaseModel: database.BaseModel{
				ID: search.ID,
			},
		})
	}
	if search.Name != nil {
		db = db.Where(group.Group{
			Name: search.Name,
		})
	}
	if search.CategoryId != uuid.Nil {
		db = db.Where(`cast(? as uuid) = any("groups"."category_ids" :: uuid[])`, search.CategoryId)
	}
	return db
}

func (u *ServerCDBRepo) InsertGroup(ctx context.Context, value *group.Group) (*group.Group, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.InsertGroup))
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

func (u *ServerCDBRepo) UpdateGroup(ctx context.Context, search *group.Search, value *group.Group) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.UpdateGroup))
	defer span.End()

	if err := applySearchGroup(u.Db, search).WithContext(ctx).Updates(value).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}

	return nil
}

func (u *ServerCDBRepo) SelectGroup(ctx context.Context, search *group.Search) (*group.Group, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.SelectCategory))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r := group.Group{}

	if err := applySelectGroup(u.Db, search).WithContext(ctx).First(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &r, nil
}

func (u *ServerCDBRepo) ListGroups(ctx context.Context, search *group.Search) ([]*group.Group, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.ListGroups))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r := make([]*group.Group, 0)

	db := u.Db
	if err := applySearchGroup(db, search).WithContext(ctx).
		Joins(`LEFT JOIN "categories" on "categories"."id" = any("groups"."category_ids" :: uuid[])`).
		Where(`"groups"."id" IN (?)`, getSubQuerySearch(db, search)).
		Select(search.Fields).Group(`"groups"."id", "groups"."name", "groups"."fee"`).Find(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return r, nil

}

func (u *ServerCDBRepo) TotalGroup(ctx context.Context, search *group.Search) (*int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.TotalGroup))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	var r int64
	if err := applySearchGroup(u.Db, search).Model(group.Group{}).
		WithContext(ctx).Count(&r).Error; err != nil {
		lib.RecordError(span, err, ctx)
		return nil, xerrors.Errorf("%w", err)
	}
	return &r, nil
}

func (u *ServerCDBRepo) CheckCategoryIdExisted(ctx context.Context, search *group.Search, categoryId uuid.UUID) (*group.Group, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.CheckCategoryIdExisted))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r := &group.Group{}

	if err := applyCheckExisted(u.Db, search).WithContext(ctx).
		Where(`cast(? as uuid) = any("groups"."category_ids" :: uuid[])`, categoryId).
		Select(`"groups".*`).First(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return r, nil

}

func getSubQuerySearch(db *gorm.DB, search *group.Search) *gorm.DB {
	if search.CategoryId == uuid.Nil {
		db = db.Table("groups")
		db = db.Select(`"groups"."id"`)
		return db
	}
	db = db.Table("groups")
	db = db.Select(`"groups"."id"`)
	db = db.Where(`cast(? as uuid) = any("groups"."category_ids" :: uuid[])`, search.CategoryId)
	return db
}

func applyCheckExisted(db *gorm.DB, search *group.Search) *gorm.DB {
	if search.ID != uuid.Nil {
		db = db.Where(`"groups"."id" != ?`, search.ID)
	}
	return db
}
