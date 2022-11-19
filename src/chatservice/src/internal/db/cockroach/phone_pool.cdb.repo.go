package cockroach

import (
	"context"

	"github.com/aqaurius6666/chatservice/src/internal/db/phone_pool"
	"github.com/aqaurius6666/chatservice/src/internal/lib"
	"github.com/aqaurius6666/chatservice/src/internal/var/c"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
	"github.com/ubgo/gormuuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func applySearchPhonePool(db *gorm.DB, search *phone_pool.Search) *gorm.DB {
	if search.ID != uuid.Nil {
		db = db.Where(phone_pool.PhonePool{
			BaseModel: database.BaseModel{
				ID: search.ID,
			},
		})
	}
	if search.PhoneNumber != nil {
		db = db.Where(phone_pool.PhonePool{
			PhoneNumber: search.PhoneNumber,
		})
	}
	if search.DefaultSearchModel.OrderBy != "" && search.DefaultSearchModel.OrderType != "" {
		db = db.Order(search.DefaultSearchModel.OrderBy + " " + search.DefaultSearchModel.OrderType)
	}

	return db
}

// Implement ListPhonePool
func (s *ServerCDBRepo) ListPhonePool(ctx context.Context, search *phone_pool.Search) ([]*phone_pool.PhonePool, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListPhonePool))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	r := make([]*phone_pool.PhonePool, 0)
	if err := applySearchPhonePool(s.Db, search).WithContext(ctx).Find(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	return r, nil
}

func (s *ServerCDBRepo) SelectPhonePool(ctx context.Context, search *phone_pool.Search) (*phone_pool.PhonePool, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.SelectPhonePool))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r := phone_pool.PhonePool{}
	if err := applySearchPhonePool(s.Db, search).WithContext(ctx).First(&r).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = phone_pool.ErrNotFound
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	return &r, nil
}

func (s *ServerCDBRepo) InsertPhonePool(ctx context.Context, value *phone_pool.PhonePool) (*phone_pool.PhonePool, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.InsertPhonePool))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := s.Db.WithContext(ctx).Create(value).Error; err != nil {
		err = xerrors.Errorf("%w", phone_pool.ErrInsertFail)
		lib.RecordError(span, err)
		return nil, err
	}
	return value, nil
}

func (s *ServerCDBRepo) ListPhonePools(ctx context.Context, search *phone_pool.Search) ([]*phone_pool.PhonePool, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListPhonePools))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	r := make([]*phone_pool.PhonePool, 0)
	if err := applySearchPhonePool(s.Db, search).WithContext(ctx).Find(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	return r, nil
}

func (s *ServerCDBRepo) UpdatePhonePool(ctx context.Context, search *phone_pool.Search, value *phone_pool.PhonePool) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UpdatePhonePool))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := applySearchPhonePool(s.Db, search).WithContext(ctx).Updates(value).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	return nil
}

func (s *ServerCDBRepo) ListAvailablePhone(ctx context.Context, phones []string) ([]*phone_pool.PhonePool, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListAvailablePhone))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	db := s.Db.WithContext(ctx).Begin()

	ret := make([]*phone_pool.PhonePool, 0)

	exprs := make([]clause.Expression, 0)
	for _, p := range phones {
		exprs = append(exprs, clause.Expr{
			SQL:  `cast(? as varchar(16)) = any(phone_number_members)`,
			Vars: []interface{}{p},
		})
	}
	var subQueryResult struct {
		PhonePoolId gormuuid.UUIDArray `gorm:"column:phone_pool_id"`
	}
	if err := db.Table("conversations").Select("phone_pool_id").Where(clause.AndConditions{
		Exprs: []clause.Expression{
			clause.Eq{
				Column: "status",
				Value:  0,
			},
			clause.OrConditions{
				Exprs: exprs,
			},
		},
	}).Find(&subQueryResult).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		db.Rollback()
		return nil, err
	}
	db = db.Where(`status = 0`)
	if len(subQueryResult.PhonePoolId) != 0 {
		db = db.Where(`id not in (?)`, subQueryResult.PhonePoolId)
	}

	if err := db.Debug().Find(&ret).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		db.Rollback()
		return nil, err
	}
	db.Commit()
	return ret, nil
}
