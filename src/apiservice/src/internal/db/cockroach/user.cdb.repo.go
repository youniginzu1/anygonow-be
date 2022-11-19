package cockroach

import (
	"context"

	"github.com/aqaurius6666/apiservice/src/internal/db/user"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

func applySearchUser(db *gorm.DB, search *user.Search) *gorm.DB {
	if search.ID != uuid.Nil {
		db = db.Where(user.User{
			BaseModel: database.BaseModel{
				ID: search.ID,
			},
		})
	}
	if search.Query != nil {
		db = db.Where(`"users"."mail" = ? OR "users"."phone" = ?`, *search.Query, *search.Query)
	}
	if search.Status != nil {
		db = db.Where(user.User{
			Status: search.Status,
		})
	}
	if search.Skip != 0 {
		db = db.Offset(search.Skip)
	}
	if search.Limit != 0 {
		db = db.Limit(search.Limit)
	}
	if search.Mail != nil {
		db = db.Where(`"users"."mail" like ?`, *search.Mail+"%")
	}
	if search.Phone != nil {
		db = db.Where(`"users"."phone" like ?`, *search.Phone+"%")
	}
	return db
}

func (u *ServerCDBRepo) SelectUser(ctx context.Context, search *user.Search) (*user.User, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.SelectUser))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r := user.User{}
	if err := applySearchUser(u.Db, search).WithContext(ctx).Joins("Contact").First(&r).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = user.ErrNotFound
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &r, nil
}

func (u *ServerCDBRepo) SelectUserProcess(ctx context.Context, search *user.Search) (*user.User, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.SelectUser))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r := user.User{}
	if err := applySearchUser(u.Db, search).WithContext(ctx).Joins("Contact").Select(search.Fields).First(&r).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = user.ErrNotFound
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &r, nil
}

func (u *ServerCDBRepo) InsertUser(ctx context.Context, value *user.User) (*user.User, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.InsertUser))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := u.Db.WithContext(ctx).Create(value).Error; err != nil {
		err = xerrors.Errorf("%w", user.ErrInsertFail)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return value, nil
}

func (u *ServerCDBRepo) ListUsers(ctx context.Context, search *user.Search) ([]*user.User, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.ListUsers))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	r := make([]*user.User, 0)
	if err := applySearchUser(u.Db, search).WithContext(ctx).Select(search.Fields).Joins("Contact").Find(&r).Error; err != nil {
		return nil, xerrors.Errorf("%w", err)
	}
	return r, nil
}

func (u *ServerCDBRepo) UpdateUser(ctx context.Context, search *user.Search, value *user.User) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.UpdateUser))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	if err := applySearchUser(u.Db, search).WithContext(ctx).Updates(value).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}

func (u *ServerCDBRepo) TotalUsers(ctx context.Context, search *user.Search) (*int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.TotalUsers))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	var r int64
	if err := applySearchUser(u.Db, search).WithContext(ctx).Model(user.User{}).Count(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &r, nil
}

func (u *ServerCDBRepo) DeleteUser(ctx context.Context, search *user.Search) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.DeleteUser))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	if err := applySearchUser(u.Db, search).WithContext(ctx).Delete(&user.User{}).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}
