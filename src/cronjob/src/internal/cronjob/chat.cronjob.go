package cronjob

import (
	"context"
	"fmt"
	"time"

	"github.com/aqaurius6666/cronjob/src/internal/lib"
	"github.com/aqaurius6666/cronjob/src/internal/model"
	"github.com/aqaurius6666/cronjob/src/internal/var/c"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/xerrors"
)

type ChatCronjob struct {
	Logger           *logrus.Logger
	Model            model.Server
	QuantityInterval CHAT_QUANTITY_INTERVAL
}

func (s *ChatCronjob) Run(ctx context.Context) {
	defer s.tearDown()
	timeCh := s.timer(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-timeCh:
			ctx2, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			s.execute(ctx2)
		}
	}
}

func (s *ChatCronjob) tearDown() {
	fmt.Println("Job chat says \"Bye bye ~\"")
}

func (s *ChatCronjob) timer(ctx context.Context) <-chan bool {
	ch := make(chan bool)
	go func() {
		ticker := time.NewTicker(time.Duration(s.QuantityInterval) * time.Minute)
		for {
			<-ticker.C
			ch <- true
		}
	}()

	return ch
}

func (s *ChatCronjob) execute(ctx context.Context) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.execute))
	defer span.End()
	startTime := time.Now()
	//

	inactiveUsers, err := s.Model.GetInactiveList(ctx)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	if len(inactiveUsers) <= 0 {
		return nil
	}
	defer func() {
		endtime := time.Now()
		s.Logger.Info("End execute chat job at ", endtime)
		span.SetAttributes(attribute.KeyValue{
			Key:   attribute.Key("job.end_time"),
			Value: attribute.StringValue(endtime.String()),
		})
	}()
	s.Logger.Info("Start execute chat job at ", startTime)
	span.SetAttributes(attribute.KeyValue{
		Key:   attribute.Key("job.start_time"),
		Value: attribute.StringValue(startTime.String()),
	})
	s.Logger.Info("Inactive users: ", inactiveUsers)
	for _, user := range inactiveUsers {
		if user.InactiveTime.Sub(startTime) > 0 {
			break
		}
		err := s.Model.TriggerSendSMS(ctx, user.UserId)
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err)
			return err
		}
	}
	return nil
}
