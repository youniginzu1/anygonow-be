package payment

import (
	"context"
	"fmt"
	"time"

	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/lib/unleash"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/apiservice/src/internal/var/e"
	"github.com/aqaurius6666/apiservice/src/pb"
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
func (s *Stripe) GetPaymentIntentAmount(ctx context.Context, id *string) (int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetPaymentIntentAmount))
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	si, err := s.client.PaymentIntents.Get(*id, &stripe.PaymentIntentParams{
		Params: stripe.Params{
			Context: ctx,
		},
	})
	if err != nil {
		lib.RecordError(span, err, ctx)
		err = xerrors.Errorf("%w", e.ErrPayment)
		return 0, err
	}
	if si.Status != stripe.PaymentIntentStatusSucceeded {
		s.logger.Error("payment intent failed")
		err = xerrors.Errorf("%w", xerrors.New("payment intent failed"))
		lib.RecordError(span, err)
	}
	return si.AmountCapturable, nil
}

func (s *Stripe) SetupIntent(ctx context.Context, cus *string) (*string, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.SetupIntent))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	si, err := s.client.SetupIntents.New(&stripe.SetupIntentParams{
		Params: stripe.Params{
			Context: ctx,
		},
		Customer:           cus,
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Usage:              &ON_SESSION,
	})
	if err != nil {
		lib.RecordError(span, err, ctx)
		err = xerrors.Errorf("%w", e.ErrPayment)
		return nil, err
	}
	return &si.ClientSecret, nil
}

func (s *Stripe) CreateCustomer(ctx context.Context, id string, email string) (string, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CreateCustomer))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	customer, err := s.client.Customers.New(&stripe.CustomerParams{
		Params: stripe.Params{
			Context: ctx,
			Metadata: map[string]string{
				"id": id,
			},
		},
		Email: &email,
	})
	if err != nil {
		lib.RecordError(span, err, ctx)
		err = xerrors.Errorf("%w", e.ErrPayment)
		return "", err
	}
	return customer.ID, nil
}

func (s *Stripe) ConfirmSetupIntent(ctx context.Context, siId string, clientSecret string) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ConfirmSetupIntent))
	defer span.End()

	s.client.SetupIntents.Confirm(siId, &stripe.SetupIntentConfirmParams{
		Params: stripe.Params{
			Context: ctx,
		},
	})
	return nil
}

func (s *Stripe) GetPublicKey(ctx context.Context) (string, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetPublicKey))
	defer span.End()

	key := unleash.GetVariant("stripe.public_key")
	return key, nil
}

func (s *Stripe) GetPaymentMethodInfo(ctx context.Context, id *string) (*pb.PaymentMethodInfo, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetPaymentMethodInfo))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	pm, err := s.client.PaymentMethods.Get(*id, &stripe.PaymentMethodParams{
		Params: stripe.Params{
			Context: ctx,
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &pb.PaymentMethodInfo{
		CardType:   string(pm.Card.Brand),
		Last4:      pm.Card.Last4,
		ExpireDate: fmt.Sprintf("%02d/%04d", pm.Card.ExpMonth, pm.Card.ExpYear),
		OwnerName:  pm.BillingDetails.Name,
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
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return (*string)(&pi.Status), nil
}
