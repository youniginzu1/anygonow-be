package cockroach

import (
	"context"
	"time"

	"github.com/aqaurius6666/authservice/src/internal/db/role"
	"github.com/aqaurius6666/authservice/src/internal/lib"
	"github.com/aqaurius6666/authservice/src/internal/var/c"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

func applySearchRole(db *gorm.DB, search *role.Search) *gorm.DB {
	if search.ID != uuid.Nil {
		db = db.Where(role.Role{
			BaseModel: database.BaseModel{
				ID: search.ID,
			},
		})
	}
	if search.Code != nil {
		db = db.Where(role.Role{
			Code: search.Code,
		})
	}
	return db
}

func (u *ServerCDBRepo) SelectRole(ctx context.Context, search *role.Search) (*role.Role, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.InsertRole))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	r := role.Role{}
	if err := applySearchRole(u.Db, search).WithContext(ctx).First(&r).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = role.ErrNotFound
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	return &r, nil
}

func (u *ServerCDBRepo) InsertRole(ctx context.Context, value *role.Role) (*role.Role, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.InsertRole))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	if err := u.Db.WithContext(ctx).Create(value).Error; err != nil {
		err = xerrors.Errorf("%w", role.ErrInsertFail)
		lib.RecordError(span, err)
		return nil, err
	}
	return value, nil
}
