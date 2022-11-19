package model

import (
	"context"

	"github.com/aqaurius6666/cronjob/src/internal/db/transaction"
	"github.com/aqaurius6666/cronjob/src/internal/lib"
	"github.com/aqaurius6666/cronjob/src/internal/var/c"
	"github.com/aqaurius6666/cronjob/src/internal/var/e"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

var (
	_ TransactionModel = (*ServerModel)(nil)
)

type TransactionModel interface {
	Pay(ctx context.Context, timeStart int64) ([]*transaction.Transaction, error)
	UpdateTransactionAfterPay(ctx context.Context, timeStart int64, bId uuid.UUID, paymentIntentId *string) error
	ChargeMoney(ctx context.Context, transaction *transaction.Transaction) (*string, error)
}

func (s *ServerModel) Pay(ctx context.Context, timeStart int64) ([]*transaction.Transaction, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.Pay))
	defer span.End()

	trans, err := s.Repo.GetTotalMoneyOfBusinessesPreviousMonth(ctx, &transaction.Search{
		DefaultSearchModel: database.DefaultSearchModel{
			Fields: []string{`"transactions"."business_id", "payments"."customer_id" as "stripe_customer_id", "payments"."payment_method_id", SUM("transactions"."fee") as "fee"`},
		},
		Now: &timeStart,
		Transaction: transaction.Transaction{
			IsPaid: lib.SafeBoolPtr(false),
			IsFree: lib.SafeBoolPtr(false),
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}

	return trans, nil
}

func (s *ServerModel) UpdateTransactionAfterPay(ctx context.Context, timeStart int64, bId uuid.UUID, paymentIntentId *string) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.Pay))
	defer span.End()

	err := s.Repo.UpdateTransaction(ctx, &transaction.Search{
		Transaction: transaction.Transaction{
			IsFree:     lib.SafeBoolPtr(false),
			BusinessId: bId,
		},
		Now: &timeStart,
	}, &transaction.Transaction{
		PaymentIntentId: paymentIntentId,
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}

	return nil
}

func (s *ServerModel) ChargeMoney(ctx context.Context, transaction *transaction.Transaction) (*string, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ChargeMoney))
	defer span.End()

	if transaction.PaymentMethodId == nil {
		err := xerrors.Errorf("%w", e.ErrMissingPaymentMethod)
		lib.RecordError(span, err)
		return nil, err
	}

	paymentIntentId, err := s.Payment.PaymentIntent(ctx, transaction.StripeCustomerId, transaction.Fee, transaction.PaymentMethodId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}

	return paymentIntentId, nil
}
