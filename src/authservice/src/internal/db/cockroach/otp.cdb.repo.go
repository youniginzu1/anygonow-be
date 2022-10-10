package cockroach

import (
	"context"
	"time"

	"github.com/aqaurius6666/authservice/src/internal/db/otp"
	"github.com/aqaurius6666/authservice/src/internal/lib"
	"github.com/aqaurius6666/authservice/src/internal/var/c"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

func applySearchOtp(db *gorm.DB, search *otp.Search) *gorm.DB {
	if search.ID != uuid.Nil {
		db = db.Where(otp.Otp{
			BaseModel: database.BaseModel{
				ID: search.ID,
			},
		})
	}
	if search.Commited != nil {
		db = db.Where(otp.Otp{
			Commited: search.Commited,
		})
	}

	return db
}

func (u *ServerCDBRepo) SelectOtp(ctx context.Context, search *otp.Search) (*otp.Otp, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.SelectOtp))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	r := otp.Otp{}
	if err := applySearchOtp(u.Db, search).WithContext(ctx).First(&r).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = otp.ErrNotFound
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	return &r, nil
}

func (u *ServerCDBRepo) InsertOtp(ctx context.Context, value *otp.Otp) (*otp.Otp, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.InsertOtp))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	if err := u.Db.WithContext(ctx).Create(value).Error; err != nil {
		err = xerrors.Errorf("%w", otp.ErrInsertFail)
		lib.RecordError(span, err)
		return nil, err
	}
	return value, nil
}

func (u *ServerCDBRepo) DeleteOTP(ctx context.Context, s *otp.Search) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.DeleteOTP))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	var r otp.Otp
	if err := applySearchOtp(u.Db, s).WithContext(ctx).Delete(&r).Error; err != nil {
		err = xerrors.Errorf("%w", otp.ErrDeleteFail)
		lib.RecordError(span, err)
		return err
	}
	return nil
}

func (u *ServerCDBRepo) UpdateOTP(ctx context.Context, search *otp.Search, value *otp.Otp) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.UpdateOTP))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	if err := applySearchOtp(u.Db, search).WithContext(ctx).Updates(value).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	return nil
}
