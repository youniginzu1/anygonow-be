package conversation

import "context"

type ConversationRepo interface {
	InsertConversation(context.Context, *Conversation) (*Conversation, error)
	SelectConversation(context.Context, *Search) (*Conversation, error)
	ListConversations(context.Context, *Search) ([]*Conversation, error)
	UpdateConversation(context.Context, *Search, *Conversation) error
	ListUnusedConversation(context.Context, *Search) ([]*Conversation, error)
	SetConversationPhonePoolNull(context.Context, *Search) error
}
