package payment

import "context"

type PaymentRepo interface {
	SelectPayment(context.Context, *Search) (*Payment, error)
	InsertPayment(context.Context, *Payment) (*Payment, error)
	UpdatePayment(context.Context, *Search, *Payment) error
	ListPayments(context.Context, *Search) ([]*Payment, error)
}
