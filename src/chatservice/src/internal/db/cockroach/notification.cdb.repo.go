package cockroach

import (
	"context"

	"github.com/aqaurius6666/chatservice/src/internal/db/notification"
	"github.com/aqaurius6666/chatservice/src/internal/lib"
	"github.com/aqaurius6666/chatservice/src/internal/var/c"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func applySearchNotification(db *gorm.DB, search *notification.Search) *gorm.DB {
	if search.UserId != uuid.Nil {
		db = db.Where(notification.Notification{
			UserId: search.UserId,
		})
	}
	return db
}

func (s *ServerCDBRepo) SelectNotification(ctx context.Context, search *notification.Search) (*notification.Notification, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.SelectNotification))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r := notification.Notification{}
	if err := applySearchNotification(s.Db, search).WithContext(ctx).Find(&r).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = notification.ErrNotFound
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	return &r, nil
}

func (s *ServerCDBRepo) UpdateNotification(ctx context.Context, search *notification.Search, value *notification.Notification) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UpdateNotification))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := applySearchNotification(s.Db, search).WithContext(ctx).Updates(value).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	return nil
}

func (s *ServerCDBRepo) InsertNotification(ctx context.Context, value *notification.Notification) (*notification.Notification, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.InsertConversation))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := s.Db.WithContext(ctx).Create(value).Error; err != nil {
		err = xerrors.Errorf("%w", notification.ErrInsertFail)
		lib.RecordError(span, err)
		return nil, err
	}
	return value, nil
}

func (s *ServerCDBRepo) UpsertNotification(ctx context.Context, value *notification.Notification) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UpsertNotification))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := s.Db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"seen"}),
	}).Create(value).Error; err != nil {
		err = xerrors.Errorf("%w", notification.ErrFirstOrCreateFail)
		lib.RecordError(span, err)
		return err
	}
	return nil
}
