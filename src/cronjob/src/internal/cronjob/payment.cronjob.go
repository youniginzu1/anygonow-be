package cronjob

import (
	"context"
	"fmt"
	"time"

	"github.com/aqaurius6666/cronjob/src/internal/db/transaction"
	"github.com/aqaurius6666/cronjob/src/internal/lib"
	"github.com/aqaurius6666/cronjob/src/internal/model"
	"golang.org/x/xerrors"

	"github.com/aqaurius6666/cronjob/src/internal/var/c"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type PaymentCronjob struct {
	Logger           *logrus.Logger
	Model            model.Server
	QuantityInterval QUANTITY_INTERVAL
	Unit             UNIT_INTERVAL
	PaymentDay       PAYMENT_DAY
}

func (s *PaymentCronjob) Run(ctx context.Context) {
	defer s.tearDown()
	if s.shouldExecute(ctx) {
		ctx2, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		s.execute(ctx2)
	}
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

func (s *PaymentCronjob) tearDown() {
	fmt.Println("Job payment says \"Bye bye ~\"")
}

func (s *PaymentCronjob) timer(ctx context.Context) <-chan bool {
	ch := make(chan bool)
	go func() {
		nextPeriod := s.getNextPeriod()
		<-time.After(time.Until(nextPeriod))
		ch <- true
		ticker := time.NewTicker(lib.GetTimeRange(int(s.QuantityInterval), string(s.Unit)))
		for {
			<-ticker.C
			ch <- true
		}
	}()

	return ch
}

func (s *PaymentCronjob) execute(ctx context.Context) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.execute))
	defer span.End()
	startTime := time.Now()
	// get detail transaction
	trans, err := s.Model.Pay(ctx, startTime.UnixMilli())
	if err != nil {
		s.Logger.Error("Cannot get transaction detail from database", time.Now())
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	if len(trans) <= 0 {
		return nil
	}
	defer func() {
		endtime := time.Now()
		s.Logger.Info("End execute job at ", endtime)
		span.SetAttributes(attribute.KeyValue{
			Key:   attribute.Key("job.end_time"),
			Value: attribute.StringValue(endtime.String()),
		})
	}()
	s.Logger.Info("Start execute job at ", startTime)
	span.SetAttributes(attribute.KeyValue{
		Key:   attribute.Key("job.start_time"),
		Value: attribute.StringValue(startTime.String()),
	})
	result := make(chan struct{}, len(trans))
	go func() {
		blocking := make(chan struct{}, c.MAX_CONCURRENCY)
		for _, v := range trans {
			blocking <- struct{}{}
			go func(v *transaction.Transaction) {
				s.chargeMoney(ctx, v, startTime.UnixMilli())
				result <- <-blocking
			}(v)
		}
	}()
	for i := 0; i < len(trans); i++ {
		<-result
	}
	return nil
}

func (s *PaymentCronjob) shouldExecute(ctx context.Context) bool {

	return false
}

func (s *PaymentCronjob) chargeMoney(ctx context.Context, tran *transaction.Transaction, timeStart int64) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.execute))
	defer span.End()
	paymentIntentId, err := s.Model.ChargeMoney(ctx, tran)
	if err != nil {
		s.Logger.Error("Cannot charge money", time.Now())
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return
	}
	err = s.Model.UpdateTransactionAfterPay(ctx, timeStart, tran.BusinessId, paymentIntentId)
	if err != nil {
		s.Logger.Error("Something wrong when update transaction status", time.Now())
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return
	}
}

func (s *PaymentCronjob) getNextPeriod() time.Time {
	return externalGetNextPeriod(string(s.PaymentDay), string(s.Unit), time.Now())
}

func externalGetNextPeriod(day string, unit string, now time.Time) time.Time {
	if unit == "WEEK" {
		daysDiff := (7 + parseWeekday(day) - now.Weekday()) % 7
		if daysDiff == 0 {
			daysDiff = 7
		}
		h, m, s := now.Clock()
		dura, _ := time.ParseDuration(fmt.Sprintf("%dh%dm%ds", h, m, s))
		nt := now.AddDate(0, 0, int(daysDiff)).Add(-dura)
		return nt
	}
	if unit == "DAY" {
		h, m, s := now.Clock()
		dura, _ := time.ParseDuration(fmt.Sprintf("%dh%dm%ds", h, m, s))
		nt := now.AddDate(0, 0, 1).Add(-dura)
		return nt
	}
	if unit == "HOUR" {
		_, m, s := now.Clock()
		dura, _ := time.ParseDuration(fmt.Sprintf("%dm%ds", m, s))
		nt := now.Add(1 * time.Hour).Add(-dura)
		return nt
	}
	return time.Time{}
}

var daysOfWeek = map[string]time.Weekday{
	"Sunday":    time.Sunday,
	"Monday":    time.Monday,
	"Tuesday":   time.Tuesday,
	"Wednesday": time.Wednesday,
	"Thursday":  time.Thursday,
	"Friday":    time.Friday,
	"Saturday":  time.Saturday,
}

func parseWeekday(v string) time.Weekday {
	if d, ok := daysOfWeek[v]; ok {
		return d
	}
	panic("day err")
}
