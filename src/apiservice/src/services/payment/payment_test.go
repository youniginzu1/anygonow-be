package payment

import (
	"context"
	"testing"

	"github.com/aqaurius6666/go-utils/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	KEY = "sk_test_51KQyIYHKYRxtUDcpsHIX2cn6KOHkSDS3Am0iKNKXesTuYkxu87qeibI5UFHh8U1LITM4SVlBJRQD6r3ChnX0C2DH00LcZU7b5j"
)

func TestCreateCustomer(t *testing.T) {
	ctx := context.Background()
	stp, err := NewStripe(ctx, logrus.New(), STRIPE_API_KEY(KEY))
	assert.Nil(t, err)
	_, err = stp.CreateCustomer(ctx, "aqaurius@gmail.com", "123")
	assert.Nil(t, err)
}

func TestSetupPaymentIntent(t *testing.T) {
	ctx := context.Background()
	stp, err := NewStripe(ctx, logrus.New(), STRIPE_API_KEY(KEY))
	assert.Nil(t, err)
	_, err = stp.SetupIntent(ctx, utils.StrPtr("cus_LAl04uGD9pxhbb"))
	assert.Nil(t, err)
}

func TestGetPublicKey(t *testing.T) {
	ctx := context.Background()
	stp, err := NewStripe(ctx, logrus.New(), STRIPE_API_KEY(KEY))
	assert.Nil(t, err)
	str, err := stp.GetPublicKey(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, str)
}
