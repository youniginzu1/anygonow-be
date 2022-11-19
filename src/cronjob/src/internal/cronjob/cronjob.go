package cronjob

import (
	"context"
	"time"

	"github.com/google/wire"
)

type QUANTITY_INTERVAL int
type UNIT_INTERVAL string
type PAYMENT_DAY string
type CHAT_QUANTITY_INTERVAL int

var JobSet = wire.NewSet(SliceProvider, wire.Struct(new(PaymentCronjob), "*"), wire.Struct(new(ChatCronjob), "*"))

var timeout = 20 * time.Second

func SliceProvider(a *PaymentCronjob, b *ChatCronjob) []Cronjob {
	return []Cronjob{
		a, b,
	}
}

type Cronjob interface {
	Run(context.Context)
}
