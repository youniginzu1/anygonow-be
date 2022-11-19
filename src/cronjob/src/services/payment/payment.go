package payment

import (
	"context"

	"github.com/google/wire"
)

var (
	Set = wire.NewSet(wire.Bind(new(Interface), new(*Stripe)), NewStripe)
)

type Interface interface {
	PaymentIntent(ctx context.Context, customerId *string, amount *int64, paymentMethodId *string) (*string, error)
}
