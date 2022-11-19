package model

import (
	"context"
	"strconv"
	"time"

	"github.com/aqaurius6666/chatservice/src/internal/db/notification"
	"github.com/aqaurius6666/chatservice/src/internal/lib"
	"github.com/aqaurius6666/chatservice/src/internal/lib/unleash"
	"github.com/aqaurius6666/chatservice/src/internal/var/c"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

type Notification interface {
	GetNotification(context.Context, interface{}) (*notification.Notification, error)
	UpsertNotification(context.Context, *notification.Notification) error

	SetUserInactiveTimeout(ctx context.Context, userId uuid.UUID) error
}

// Imeplement SetUserInactiveTimeout
func (s *ServerModel) SetUserInactiveTimeout(ctx context.Context, userId uuid.UUID) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.SetUserInactiveTimeout))
	defer span.End()
	err := s.Redis.AddMemeberSortedSet(ctx, c.INACTIVE_SET_KEY, userId.String(), float64(time.Now().Add(s.getActiveTimeout()).UnixMilli()))
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	return nil
}

func (s *ServerModel) getActiveTimeout() time.Duration {
	v := unleash.GetVariant("chatservice.active-timeout")
	d, err := strconv.Atoi(v)
	if err != nil {
		s.Logger.Error(err)
		return 5 * time.Second
	}
	return time.Duration(d) * time.Second
}

func (s *ServerModel) GetNotification(ctx context.Context, userId interface{}) (*notification.Notification, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetNotification))
	defer span.End()

	uid, err := lib.ToUUID(userId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}

	noti, err := s.Repo.SelectNotification(ctx, &notification.Search{
		Notification: notification.Notification{
			UserId: uid,
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	err = s.SetUserInactiveTimeout(ctx, noti.UserId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}

	return noti, nil
}

func (s *ServerModel) UpsertNotification(ctx context.Context, noti *notification.Notification) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetNotification))
	defer span.End()

	err := s.Repo.UpsertNotification(ctx, noti)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}

	err = s.Redis.AddMemeberSortedSet(ctx, c.INACTIVE_SET_KEY, noti.UserId.String(), float64(time.Now().Add(s.getActiveTimeout()).UnixMilli()))
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}

	return nil
}
