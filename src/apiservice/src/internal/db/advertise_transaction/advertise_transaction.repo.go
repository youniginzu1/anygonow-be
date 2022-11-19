package advertise_transaction

import "context"

type AdvertiseTransactionRepo interface {
	SelectAdvertiseTransaction(*Search) (*AdvertiseTransaction, error)
	InsertAdvertiseTransaction(context.Context, *AdvertiseTransaction) (*AdvertiseTransaction, error)
	UpdateAdvertiseTransaction(*Search, *AdvertiseTransaction) error
}
