package model

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aqaurius6666/go-utils/utils"
	"github.com/aqaurius6666/mailservice/src/internal/db/device"
	"github.com/aqaurius6666/mailservice/src/internal/lib"
	"github.com/aqaurius6666/mailservice/src/internal/var/c"
	"github.com/aqaurius6666/mailservice/src/internal/var/e"
	"github.com/aqaurius6666/mailservice/src/pb/mailpb"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

type NotificationModel interface {
	SubscribeNotification(ctx context.Context, req *mailpb.SubscribeNotificationRequest) (*mailpb.SubscribeNotificationResponse, error)
	UnsubscribeNotification(ctx context.Context, req *mailpb.UnsubscribeNotificationRequest) (*mailpb.UnsubscribeNotificationResponse, error)
	SendNotification(ctx context.Context, req *mailpb.SendNotificationRequest) (*mailpb.SendNotificationResponse, error)
}

func (s *ServerModel) SubscribeNotification(ctx context.Context, req *mailpb.SubscribeNotificationRequest) (*mailpb.SubscribeNotificationResponse, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.SubscribeNotification))
	defer span.End()

	var (
		dev *device.Device
		err error
	)
	dev, err = s.Repo.SelectDevice(ctx, &device.Search{
		Device: device.Device{
			DeviceId: &req.DeviceId,
			UserId:   &req.UserId,
		},
	})
	if err == nil {
		if *dev.Active == device.ACTIVE {
			return &mailpb.SubscribeNotificationResponse{}, nil
		}
		err = s.Repo.UpdateDevice(ctx, &device.Search{
			Device: device.Device{
				DeviceId: &req.DeviceId,
				UserId:   &req.UserId,
				Active:   utils.IntPtr(device.INACTIVE),
			},
		}, &device.Device{
			Active: utils.IntPtr(device.ACTIVE),
		})
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err)
			return nil, err
		}
		return &mailpb.SubscribeNotificationResponse{}, nil
	}

	if !errors.Is(err, device.ErrNotFound) {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}

	dev, err = s.Repo.InsertDevice(ctx, &device.Device{
		DeviceId: &req.DeviceId,
		UserId:   &req.UserId,
		Active:   utils.IntPtr(device.ACTIVE),
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	s.Logger.Debug("register new device %s, userId = %s", *dev.DeviceId, *dev.UserId)
	return &mailpb.SubscribeNotificationResponse{}, nil
}

func (s *ServerModel) UnsubscribeNotification(ctx context.Context, req *mailpb.UnsubscribeNotificationRequest) (*mailpb.UnsubscribeNotificationResponse, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UnsubscribeNotification))
	defer span.End()

	err := s.Repo.UpdateDevice(ctx, &device.Search{
		Device: device.Device{
			DeviceId: &req.DeviceId,
			UserId:   &req.UserId,
			Active:   utils.IntPtr(device.ACTIVE),
		},
	}, &device.Device{
		Active: utils.IntPtr(device.INACTIVE),
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	return &mailpb.UnsubscribeNotificationResponse{}, nil
}

func (s *ServerModel) SendNotification(ctx context.Context, req *mailpb.SendNotificationRequest) (*mailpb.SendNotificationResponse, error) {
	go func() {
		ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(context.Background(), lib.GetFunctionName(s.SendNotification))
		defer span.End()

		_, err := uuid.Parse(req.To)
		if err != nil {
			err = xerrors.Errorf("%w", e.ErrIdInvalidFormat)
			lib.RecordError(span, err)
			s.Logger.Error(err)
			return
		}
		devs, err := s.Repo.ListDevices(ctx, &device.Search{
			Device: device.Device{
				UserId: &req.To,
				Active: utils.IntPtr(device.ACTIVE),
			},
		})
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err)
			s.Logger.Error(err)
			return
		}
		var message map[string]string
		err = json.Unmarshal([]byte(req.Message), &message)
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err)
			s.Logger.Error(err)
			return
		}
		validDevices := make([]string, 0)
		for _, v := range devs {
			if err := s.Fcm.CheckToken(ctx, *v.DeviceId); err == nil {
				validDevices = append(validDevices, *v.DeviceId)
				continue
			}
			err := s.Repo.DeleteDevice(ctx, &device.Search{
				Device: *v,
			}, v)
			if err != nil {
				err = xerrors.Errorf("%w", err)
				lib.RecordError(span, err)
				s.Logger.Error(err)
			}
		}
		err = s.Fcm.SendTo(ctx, validDevices, req.Title, req.Body, message)
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err)
			s.Logger.Error(err)
			return
		}
	}()
	return &mailpb.SendNotificationResponse{}, nil
}
