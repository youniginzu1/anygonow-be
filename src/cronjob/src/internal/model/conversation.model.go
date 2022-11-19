package model

import (
	"context"
	"time"

	"github.com/aqaurius6666/cronjob/src/internal/lib"
	"github.com/aqaurius6666/cronjob/src/internal/var/c"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

type ConversationModel interface {
	TriggerSendSMS(ctx context.Context, userId string) error
	GetInactiveList(ctx context.Context) ([]struct {
		UserId       string
		InactiveTime time.Time
	}, error)
}

func (s *ServerModel) TriggerSendSMS(ctx context.Context, userId string) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.TriggerSendSMS))
	defer span.End()

	err := s.Chat.TriggerSendSMS(ctx, userId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	return nil
}

// Implement GetInactiveList
func (s *ServerModel) GetInactiveList(ctx context.Context) ([]struct {
	UserId       string
	InactiveTime time.Time
}, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetInactiveList))
	defer span.End()

	list, err := s.Redis.GetSortedSet(ctx, c.INACTIVE_SET_KEY)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	ret := make([]struct {
		UserId       string
		InactiveTime time.Time
	}, len(list))
	for i, e := range list {
		ret[i] = struct {
			UserId       string
			InactiveTime time.Time
		}{
			UserId:       e.Member.(string),
			InactiveTime: time.UnixMilli(int64(e.Score)),
		}
	}
	return ret, nil
}
