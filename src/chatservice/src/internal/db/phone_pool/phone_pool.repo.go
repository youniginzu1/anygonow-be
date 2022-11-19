package phone_pool

import (
	"context"
)

type PhonePoolRepo interface {
	ListPhonePool(ctx context.Context, search *Search) ([]*PhonePool, error)
	InsertPhonePool(context.Context, *PhonePool) (*PhonePool, error)
	SelectPhonePool(context.Context, *Search) (*PhonePool, error)
	UpdatePhonePool(context.Context, *Search, *PhonePool) error

	ListAvailablePhone(context.Context, []string) ([]*PhonePool, error)
	
}
