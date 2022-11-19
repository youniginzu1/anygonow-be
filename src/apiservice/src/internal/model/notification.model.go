package model

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

type NotificationModel interface {
	SubscribeNotification(context.Context, string, string) error
	UnsubscribeNotification(context.Context, string, string) error
	SendRequestedNotification(ctx context.Context, handymanId string, customerName string, categoryName string, zipcode string) error
	SendCancelNotification(ctx context.Context, handymanId string, customerName string, categoryName string, zipcode string) error
	SendFeeNotification(ctx context.Context, handymanId string, fee float32) error

	SendConnectNotification(ctx context.Context, customerId string, businessName string, conversationId string) error
	SendCompleteNotification(ctx context.Context, customerId string, businessName string, businessId string) error
	SendRejectNotification(ctx context.Context, customerId string, businessName string, businessId string) error
}

func (s *ServerModel) SendConnectNotification(ctx context.Context, customerId string, businessName string, conversationId string) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.SendConnectNotification))
	defer span.End()
	title := "Your request"
	body := fmt.Sprintf("%s has accepted your request", businessName)

	nid := uuid.New()
	message := map[string]string{
		"id":             nid.String(),
		"seq":            fmt.Sprint(nid.ClockSequence()),
		"type":           c.CONNECT_CUSTOMER_NOTIFICATION,
		"conversationId": conversationId,
	}
	messageByte, err := json.Marshal(message)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	err = s.Mail.SendNotification(ctx, customerId, title, body, string(messageByte))
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	return nil
}

func (s *ServerModel) SendCompleteNotification(ctx context.Context, customerId string, businessName string, businessId string) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.SendCompleteNotification))
	defer span.End()
	title := "Your request"
	body := fmt.Sprintf("Your request with %s is complete", businessName)

	nid := uuid.New()
	message := map[string]string{
		"id":         nid.String(),
		"seq":        fmt.Sprint(nid.ClockSequence()),
		"type":       c.COMPLETE_CUSTOMER_NOTIFICATION,
		"businessId": businessId,
	}
	messageByte, err := json.Marshal(message)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	err = s.Mail.SendNotification(ctx, customerId, title, body, string(messageByte))
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	return nil
}

func (s *ServerModel) SendRejectNotification(ctx context.Context, customerId string, businessName string, businessId string) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.SendRejectNotification))
	defer span.End()
	title := "Your request"
	body := fmt.Sprintf("%s has rejected your request", businessName)

	nid := uuid.New()
	message := map[string]string{
		"id":         nid.String(),
		"seq":        fmt.Sprint(nid.ClockSequence()),
		"type":       c.REJECT_CUSTOMER_NOTIFICATION,
		"businessId": businessId,
	}
	messageByte, err := json.Marshal(message)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	err = s.Mail.SendNotification(ctx, customerId, title, body, string(messageByte))
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	return nil
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

func (s *ServerModel) SendRequestedNotification(ctx context.Context, handymanId string, customerName string, categoryName string, zipcode string) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.SendRequestedNotification))
	defer span.End()
	s.Logger.Info(handymanId, customerName, categoryName, zipcode)
	title := "You have a new request"
	body := fmt.Sprintf("%s sent a request for %s service at zipcode %s", customerName, categoryName, zipcode)

	nid := uuid.New()
	message := map[string]string{
		"id":   nid.String(),
		"seq":  fmt.Sprint(nid.ClockSequence()),
		"type": c.REQUEST_HANDYMAN_NOTIFICATION,
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

func (s *ServerModel) SendCancelNotification(ctx context.Context, handymanId string, customerName string, categoryName string, zipcode string) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.SendCancelNotification))
	defer span.End()
	s.Logger.Info(handymanId, customerName, categoryName, zipcode)
	title := "Request rejected"
	body := fmt.Sprintf("%s rejected request for %s service at zipcode %s", customerName, categoryName, zipcode)
	nid := uuid.New()
	message := map[string]string{
		"id":   nid.String(),
		"seq":  fmt.Sprint(nid.ClockSequence()),
		"type": c.CANCEL_HANDYMAN_NOTIFICATION,
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

func (s *ServerModel) SubscribeNotification(ctx context.Context, userId string, deviceId string) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.SubscribeNotification))
	defer span.End()

	err := s.Mail.SubscribeNotification(ctx, userId, deviceId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	return nil
}

func (s *ServerModel) UnsubscribeNotification(ctx context.Context, userId string, deviceId string) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UnsubscribeNotification))
	defer span.End()

	err := s.Mail.UnsubscribeNotification(ctx, userId, deviceId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	return nil
}
