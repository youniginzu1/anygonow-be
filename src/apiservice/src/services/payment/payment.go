package payment

import (
	"context"

	"github.com/aqaurius6666/apiservice/src/pb"
	"github.com/google/wire"
)

var (
	Set = wire.NewSet(wire.Bind(new(Interface), new(*Stripe)), NewStripe)
)

type Interface interface {
	// Checkout(context.Context) error
	// create customer_id in payment
	CreateCustomer(context.Context, string, string) (string, error)
	SetupIntent(context.Context, *string) (*string, error)
	GetPublicKey(context.Context) (string, error)
	GetPaymentMethodInfo(ctx context.Context, id *string) (*pb.PaymentMethodInfo, error)
	GetPaymentIntentAmount(ctx context.Context, id *string) (int64, error)
	PaymentIntent(ctx context.Context, customerId *string, amount *int64, paymentMethodId *string) (*string, error)
	// ConfirmSetupIntent(ctx context.Context, clientSecret *string) error
}
