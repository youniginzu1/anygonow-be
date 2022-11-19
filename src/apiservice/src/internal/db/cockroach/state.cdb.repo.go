package cockroach

import (
	"context"

	"github.com/aqaurius6666/apiservice/src/internal/db/state"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

func applySearchState(db *gorm.DB, search *state.Search) *gorm.DB {
	if search.ID != uuid.Nil {
		db = db.Where(state.State{
			BaseModel: database.BaseModel{
				ID: search.ID,
			},
		})
	}
	return db
}

func (u *ServerCDBRepo) SelectState(ctx context.Context, search *state.Search) (*state.State, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.SelectState))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r := state.State{}
	if err := applySearchState(u.Db, search).WithContext(ctx).First(&r).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = state.ErrNotFound
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &r, nil
}

func (u *ServerCDBRepo) InsertState(ctx context.Context, value *state.State) (*state.State, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.InsertState))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := u.Db.WithContext(ctx).Create(value).Error; err != nil {
		err = xerrors.Errorf("%w", state.ErrInsertFail)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return value, nil
}

func (u *ServerCDBRepo) ListStates(ctx context.Context, search *state.Search) ([]*state.State, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.ListStates))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	r := make([]*state.State, 0)
	if err := applySearchState(u.Db, search).WithContext(ctx).Find(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return r, nil
}
