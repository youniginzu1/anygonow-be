package chat

import "context"

type ChatRepo interface {
	InsertChat(context.Context, *Chat) (*Chat, error)
	SelectChat(context.Context, *Search) (*Chat, error)
	ListChats(context.Context, *Search) ([]*Chat, error)
	UpdateChat(context.Context, *Search, *Chat) error
}
