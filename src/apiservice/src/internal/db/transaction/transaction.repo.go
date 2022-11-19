package transaction

import "context"

type TransactionRepo interface {
	ListTransactions(context.Context, *Search) ([]*Transaction, error)
	TotalTransaction(context.Context, *Search) (*int64, error)
	TotalFee(context.Context, *Search) (*int64, error)
	InsertTransaction(context.Context, *Transaction) error
	UpdateTransaction(context.Context, *Search, *Transaction) error
}
