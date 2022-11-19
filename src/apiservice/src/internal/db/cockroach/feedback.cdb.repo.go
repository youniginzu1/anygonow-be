package cockroach

import (
	"context"

	"github.com/aqaurius6666/apiservice/src/internal/db/feedback"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

func applySearchFeedback(db *gorm.DB, search *feedback.Search) *gorm.DB {
	if search.ID != uuid.Nil {
		db = db.Where(feedback.Feedback{
			BaseModel: database.BaseModel{
				ID: search.ID,
			},
		})
	}
	if search.BusinessId != uuid.Nil {
		db = db.Where(feedback.Feedback{
			BusinessId: search.BusinessId,
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

func (u *ServerCDBRepo) SelectFeedback(ctx context.Context, search *feedback.Search) (*feedback.Feedback, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.SelectFeedback))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r := feedback.Feedback{}
	if err := applySearchFeedback(u.Db, search).WithContext(ctx).First(&r).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = feedback.ErrNotFound
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &r, nil
}

func (u *ServerCDBRepo) InsertFeedback(ctx context.Context, value *feedback.Feedback) (*feedback.Feedback, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.InsertFeedback))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := u.Db.WithContext(ctx).Create(value).Error; err != nil {
		err = xerrors.Errorf("%w", feedback.ErrInsertFail)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return value, nil
}

func (u *ServerCDBRepo) ListFeedbacks(ctx context.Context, search *feedback.Search) ([]*feedback.Feedback, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.ListFeedbacks))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	r := make([]*feedback.Feedback, 0)

	if err := applySearchFeedback(u.Db, search).WithContext(ctx).Find(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return r, nil
}

func (u *ServerCDBRepo) ListFeedbacksByBusiness(ctx context.Context, search *feedback.Search) ([]*feedback.Feedback, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.ListFeedbacksByBusiness))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r := make([]*feedback.Feedback, 0)
	if err := applySearchFeedback(u.Db, search).WithContext(ctx).Select(search.Fields).
		Joins(`LEFT JOIN "services" ON "services"."id" = "feedbacks"."service_id"`).
		Joins(`LEFT JOIN "orders" ON "orders"."id" = "feedbacks"."order_id"`).
		Joins(`LEFT JOIN "users" ON "users"."id" = "orders"."customer_id"`).
		Joins(`LEFT JOIN "categories" ON "categories"."id" = "services"."category_id"`).
		Order(`"feedbacks"."updated_at" DESC`).
		Find(&r).
		Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return r, nil
}

func (u *ServerCDBRepo) TotalFeedbacksByBusiness(ctx context.Context, search *feedback.Search) (*int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.TotalFeedbacksByBusiness))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var r int64
	if err := applySearchFeedback(u.Db, search).WithContext(ctx).Model(feedback.Feedback{}).Count(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &r, nil
}

func (u *ServerCDBRepo) UpdateFeedback(ctx context.Context, search *feedback.Search, value *feedback.Feedback) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.UpdateFeedback))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := applySearchFeedback(u.Db, search).WithContext(ctx).Updates(value).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}

func (u *ServerCDBRepo) SelectRatingByBussiness(ctx context.Context, search *feedback.Search) (*feedback.Feedback, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.SelectRatingByBussiness))
	defer span.End()

	if search.BusinessId == uuid.Nil {
		err := xerrors.Errorf("%w", feedback.ErrNotFound)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	v := feedback.Feedback{}
	if err := u.Db.Group("business_id").WithContext(ctx).
		Select("sum(rate) / count(rate) as rate, business_id, count(*) as count").Where(feedback.Feedback{
		BusinessId: search.BusinessId,
	}).Find(&v).Error; err != nil {
		err = xerrors.Errorf("%w", feedback.ErrNotFound)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &v, nil
}

func (u *ServerCDBRepo) ListRatingByBussiness(ctx context.Context, search *feedback.Search) ([]*feedback.Feedback, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.ListRatingByBussiness))
	defer span.End()

	if search.BusinessId == uuid.Nil {
		err := xerrors.Errorf("%w", feedback.ErrNotFound)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	v := make([]*feedback.Feedback, 0)

	if err := applySearchFeedback(u.Db, search).
		Order("rate ASC").
		Group("rate").WithContext(ctx).
		Select("rate, COUNT(*) AS review").
		Find(&v).Error; err != nil {
		err = xerrors.Errorf("%w", feedback.ErrNotFound)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return v, nil
}
