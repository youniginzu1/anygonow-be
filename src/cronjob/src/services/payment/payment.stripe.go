package payment

import (
	"context"

	"github.com/aqaurius6666/cronjob/src/internal/lib"
	"github.com/aqaurius6666/cronjob/src/internal/var/c"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

var (
	_ Interface = (*Stripe)(nil)
)

type STRIPE_API_KEY string

var (
	ON_SESSION = "on_session"
	CURRENCY   = "usd"
)

type Stripe struct {
	ctx    context.Context
	logger *logrus.Logger
	client *client.API
}

func NewStripe(ctx context.Context, logger *logrus.Logger, key STRIPE_API_KEY) (*Stripe, error) {
	b := stripe.GetBackendWithConfig(stripe.APIBackend, &stripe.BackendConfig{
		LeveledLogger: logger,
		LogLevel:      0,
	})
	c := client.New(string(key), &stripe.Backends{
		API: b,
	})
	return &Stripe{
		ctx:    ctx,
		logger: logger,
		client: c,
	}, nil
}

func (s *Stripe) PaymentIntent(ctx context.Context, customerId *string, amount *int64, paymentMethodId *string) (*string, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.PaymentIntent))
	defer span.End()

	pi, err := s.client.PaymentIntents.New(&stripe.PaymentIntentParams{
		Params: stripe.Params{
			Context: ctx,
		},
		Amount:        stripe.Int64(*amount),
		Customer:      stripe.String(*customerId),
		PaymentMethod: stripe.String(*paymentMethodId),
		Currency:      stripe.String(string(CURRENCY)),
		Confirm:       stripe.Bool(true),
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	return &pi.ID, nil
}
