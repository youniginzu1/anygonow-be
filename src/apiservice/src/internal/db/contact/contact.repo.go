package contact

import "context"

type ContactRepo interface {
	SelectContact(context.Context, *Search) (*Contact, error)
	DeleteContact(context.Context, *Search) error
	InsertContact(context.Context, *Contact) (*Contact, error)
	UpdateContact(context.Context, *Search, *Contact) error
}
