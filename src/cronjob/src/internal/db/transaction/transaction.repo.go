package transaction

import "context"

type TransactionRepo interface {
	GetTotalMoneyOfBusinessesPreviousMonth(ctx context.Context, search *Search) ([]*Transaction, error)
	UpdateTransaction(ctx context.Context, search *Search, value *Transaction) error
}
