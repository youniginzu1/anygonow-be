package fcm

import (
	"context"
	"encoding/json"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/aqaurius6666/mailservice/src/internal/lib"
	"github.com/aqaurius6666/mailservice/src/internal/var/c"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
	"google.golang.org/api/option"
)

type FB_PRIVATE_KEY string

type Service interface {
	SendTo(ctx context.Context, to []string, title, body string, message interface{}) error
	CheckToken(ctx context.Context, token string) error
}

var AppName = "AnyGoNow"
var AndroidChannelId = "anygonow"
var Set = wire.NewSet(NewService)

type ServiceImpl struct {
	App    *firebase.App
	Logger *logrus.Logger
}

func NewService(logger *logrus.Logger, key FB_PRIVATE_KEY) (Service, error) {

	var jsonData map[string]string
	err := json.Unmarshal([]byte(ACCOUNT), &jsonData)
	if err != nil {
		return nil, err
	}
	jsonData["private_key"] = string(strings.ReplaceAll(string(key), "\\n", "\n"))
	credential, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}
	opt := option.WithCredentialsJSON(credential)
	config := &firebase.Config{
		ProjectID: jsonData["project_id"],
	}
	app, err := firebase.NewApp(context.TODO(), config, opt)
	if err != nil {
		return nil, err
	}
	return &ServiceImpl{
		App:    app,
		Logger: logger,
	}, nil
}
func (s *ServiceImpl) SendTo(ctx context.Context, to []string, title, body string, message interface{}) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.SendTo))
	defer span.End()

	c, err := s.App.Messaging(ctx)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	_, err = c.SendMulticast(ctx, &messaging.MulticastMessage{
		Tokens: to,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Android: &messaging.AndroidConfig{
			Notification: &messaging.AndroidNotification{
				ChannelID: AndroidChannelId,
			},
		},
		Data: message.(map[string]string),
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	return nil
}

func (s *ServiceImpl) CheckToken(ctx context.Context, token string) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CheckToken))
	defer span.End()

	c, err := s.App.Messaging(ctx)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	_, err = c.SendDryRun(ctx, &messaging.Message{
		Token: token,
		Data: map[string]string{
			"tag": "testing",
		},
		Notification: &messaging.Notification{
			Title: AppName,
			Body:  "testing",
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	return nil
}
