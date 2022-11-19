package order

import "context"

type OrderRepo interface {
	SelectOrder(context.Context, *Search) (*Order, error)
	InsertOrder(context.Context, *Order) (*Order, error)
	UpdateOrder(context.Context, *Search, *Order) error
	TotalOrder(context.Context, *Search) (*int64, error)
	ListOrders(context.Context, *Search) ([]*Order, error)
	UpdateOrdersCancelledStatusByUser(context.Context, *Search, *Order) error
	ListProjects(context.Context, *Search) ([]*Order, error)
	TotalProjects(context.Context, *Search) (*int64, error)
	CancelProject(context.Context, *Search, *Order) error
	UpdateOrderStatusIfExpireTime(ctx context.Context) error
	ListBusinessesAlreadyOrdered(context.Context, *Search) ([]*Order, error)
}
