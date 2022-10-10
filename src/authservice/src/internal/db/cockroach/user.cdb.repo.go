package cockroach

import (
	"context"
	"time"

	"github.com/aqaurius6666/authservice/src/internal/db/user"
	"github.com/aqaurius6666/authservice/src/internal/lib"
	"github.com/aqaurius6666/authservice/src/internal/var/c"
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
	if search.PublicKey != nil {
		db = db.Where(user.User{
			PublicKey: search.PublicKey,
		})
	}
	if search.Username != nil {
		db = db.Where(user.User{
			Username: search.Username,
		})
	}
	if search.Phone != nil {
		db = db.Where(user.User{
			Phone: search.Phone,
		})
	}
	if search.Mail != nil {
		db = db.Where(user.User{
			Mail: search.Mail,
		})
	}
	return db
}

func (u *ServerCDBRepo) SelectUser(ctx context.Context, search *user.Search) (*user.User, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.SelectUser))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	r := user.User{}
	if err := applySearchUser(u.Db, search).WithContext(ctx).Joins("Role").First(&r).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = user.ErrNotFound
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	return &r, nil
}

func (u *ServerCDBRepo) InsertUser(ctx context.Context, value *user.User) (*user.User, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.InsertUser))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	if err := u.Db.WithContext(ctx).Create(value).Error; err != nil {
		err = xerrors.Errorf("%w", user.ErrInsertFail)
		lib.RecordError(span, err)
		return nil, err
	}
	return value, nil
}
func (u *ServerCDBRepo) DeleteUser(ctx context.Context, search *user.Search) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.DeleteUser))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	if err := applySearchUser(u.Db, search).WithContext(ctx).Delete(user.User{}).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	return nil
}
func (u *ServerCDBRepo) ListUsers(ctx context.Context, search *user.Search) ([]*user.User, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.ListUsers))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	r := make([]*user.User, 0)
	if err := applySearchUser(u.Db, search).WithContext(ctx).Find(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	return r, nil
}

func (u *ServerCDBRepo) UpdateUser(ctx context.Context, search *user.Search, value *user.User) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.UpdateUser))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	if err := applySearchUser(u.Db, search).WithContext(ctx).Updates(value).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	return nil
}
