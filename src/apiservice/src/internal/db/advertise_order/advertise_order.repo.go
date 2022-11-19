package advertise_order

import (
	"context"
)

type AdvertiseOrderRepo interface {
	SelectAdvertiseOrder(*Search) (*AdvertiseOrder, error)
	InsertAdvertiseOrder(context.Context, *AdvertiseOrder) (*AdvertiseOrder, error)
	UpdateAdvertiseOrder(*Search, *AdvertiseOrder) error
	TotalAdvertiseOrder(context.Context, *Search) (*int64, error)
	ListAdvertiseOrders(context.Context, *Search) ([]*AdvertiseOrder, error)
	GetTotalOrderForBuyValidate(context.Context, *Search) (*int64, error)
	TotalFeeAdvertise(context.Context, *Search) (*float64, error)
}
