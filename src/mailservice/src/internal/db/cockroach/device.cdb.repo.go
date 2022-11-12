package cockroach

import (
	"context"

	"github.com/aqaurius6666/go-utils/database"
	"github.com/aqaurius6666/mailservice/src/internal/db/device"
	"github.com/aqaurius6666/mailservice/src/internal/lib"
	"github.com/aqaurius6666/mailservice/src/internal/var/c"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

func applySearchDevice(db *gorm.DB, search *device.Search) *gorm.DB {
	if search.ID != uuid.Nil {
		db = db.Where(device.Device{
			BaseModel: database.BaseModel{
				ID: search.ID,
			},
		})
	}
	if search.Active != nil {
		db = db.Where(device.Device{
			Active: search.Active,
		})
	}
	if search.UserId != nil {
		db = db.Where(device.Device{
			UserId: search.UserId,
		})
	}
	if search.DeviceId != nil {
		db = db.Where(device.Device{
			DeviceId: search.DeviceId,
		})
	}
	return db
}

func (s *ServerCDBRepo) SelectDevice(ctx context.Context, search *device.Search) (*device.Device, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.SelectDevice))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r := device.Device{}
	if err := applySearchDevice(s.Db, search).WithContext(ctx).First(&r).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = device.ErrNotFound
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	return &r, nil
}

func (s *ServerCDBRepo) InsertDevice(ctx context.Context, value *device.Device) (*device.Device, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.InsertDevice))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := s.Db.WithContext(ctx).Create(value).Error; err != nil {
		err = xerrors.Errorf("%w", device.ErrInsertFail)
		lib.RecordError(span, err)
		return nil, err
	}
	return value, nil
}

func (s *ServerCDBRepo) ListDevices(ctx context.Context, search *device.Search) ([]*device.Device, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListDevices))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	r := make([]*device.Device, 0)
	if err := applySearchDevice(s.Db, search).WithContext(ctx).Find(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	return r, nil
}

func (s *ServerCDBRepo) UpdateDevice(ctx context.Context, search *device.Search, value *device.Device) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UpdateDevice))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := applySearchDevice(s.Db, search).WithContext(ctx).Updates(value).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	return nil
}

func (s *ServerCDBRepo) DeleteDevice(ctx context.Context, search *device.Search, v *device.Device) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.DeleteDevice))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := applySearchDevice(s.Db, search).WithContext(ctx).Delete(v).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	return nil
}
