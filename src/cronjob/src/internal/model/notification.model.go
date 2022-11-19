package model

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aqaurius6666/cronjob/src/internal/lib"
	"github.com/aqaurius6666/cronjob/src/internal/var/c"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

type NotificationModel interface {
	SendFeeNotification(ctx context.Context, handymanId string, fee float32) error
}

func (s *ServerModel) SendFeeNotification(ctx context.Context, handymanId string, fee float32) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.SendFeeNotification))
	defer span.End()
	title := "AnyGoNow"
	body := fmt.Sprintf("You have paid $%.2f for the system fee", fee)

	nid := uuid.New()
	message := map[string]string{
		"id":   nid.String(),
		"seq":  fmt.Sprint(nid.ClockSequence()),
		"type": c.FEE_HANDYMAN_NOTIFICATION,
	}
	messageByte, err := json.Marshal(message)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	err = s.Mail.SendNotification(ctx, handymanId, title, body, string(messageByte))
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	return nil
}
