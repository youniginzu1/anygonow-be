package model

import (
	"context"

	"github.com/aqaurius6666/chatservice/src/internal/db/chat"
	"github.com/aqaurius6666/chatservice/src/internal/db/conversation"
	"github.com/aqaurius6666/chatservice/src/internal/db/phone_pool"
	"github.com/aqaurius6666/chatservice/src/internal/lib"
	"github.com/aqaurius6666/chatservice/src/internal/var/c"
	"github.com/aqaurius6666/chatservice/src/pb/chatpb"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/aqaurius6666/go-utils/utils"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

type HTTP interface {
	FetchAllChats(context.Context, *chat.Search) ([]*chat.Chat, error)
	GetConversationById(context.Context, interface{}) (*conversation.Conversation, error)
	GetConversationWithMemberId(ctx context.Context, id interface{}, memberId interface{}) (*conversation.Conversation, error)
	FetchChatsByConversationId(context.Context, interface{}) ([]*chat.Chat, error)
	ConvertChatToProto(*chat.Chat) *chatpb.Chat
	ConvertChatToProtos([]*chat.Chat) []*chatpb.Chat
	CreateChat(ctx context.Context, id interface{}, userId interface{}, payload *string, shouldSeen *bool) error
	ListConversations(context.Context, *conversation.Search) ([]*conversation.Conversation, error)
	ConvertConversationToProtos(c []*conversation.Conversation) []*chatpb.Conversation
	ConvertConversationToProto(c *conversation.Conversation) *chatpb.Conversation

	SeenChat(ctx context.Context, convId uuid.UUID, receiverId uuid.UUID) error
	GetPhonePoolById(ctx context.Context, id interface{}, status *int32) (*phone_pool.PhonePool, error)
	ListPhonePool(ctx context.Context, search *phone_pool.Search) ([]*phone_pool.PhonePool, error)
}

// Implement ListPhonePool
func (s *ServerModel) ListPhonePool(ctx context.Context, search *phone_pool.Search) ([]*phone_pool.PhonePool, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListPhonePool))
	defer span.End()
	phones, err := s.Repo.ListPhonePool(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	return phones, nil
}

func (s *ServerModel) SeenChat(ctx context.Context, convId uuid.UUID, receiverId uuid.UUID) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.SeenChat))
	defer span.End()

	err := s.Repo.UpdateChat(ctx, &chat.Search{
		ReceiverId: receiverId,
		Chat: chat.Chat{
			ConversationId: convId,
		},
	}, &chat.Chat{
		Seen: utils.BoolPtr(true),
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}

	return nil
}

func (s *ServerModel) ConvertChatToProtos(c []*chat.Chat) []*chatpb.Chat {
	cpb := make([]*chatpb.Chat, 0)
	for _, cc := range c {
		cpb = append(cpb, s.ConvertChatToProto(cc))
	}
	return cpb
}

func (s *ServerModel) ConvertChatToProto(c *chat.Chat) *chatpb.Chat {
	cpb := new(chatpb.Chat)
	if c.CreatedAt != 0 {
		cpb.Timestamp = c.CreatedAt
	}
	if c.Payload != nil {
		cpb.Payload = *c.Payload
	}
	if c.SenderId != uuid.Nil {
		cpb.Sender = c.SenderId.String()
	}
	return cpb
}

func (s *ServerModel) ConvertConversationToProtos(c []*conversation.Conversation) []*chatpb.Conversation {
	cpb := make([]*chatpb.Conversation, 0)
	for _, cc := range c {
		cpb = append(cpb, s.ConvertConversationToProto(cc))
	}
	return cpb
}

func (s *ServerModel) ConvertConversationToProto(c *conversation.Conversation) *chatpb.Conversation {
	cpb := new(chatpb.Conversation)
	if c.ID != uuid.Nil {
		cpb.Id = c.ID.String()
	}
	if c.LastChat != nil {
		cpb.LastChat = *c.LastChat
	}
	if c.Members != nil {
		members := make([]string, 0)
		for _, m := range c.Members {
			members = append(members, m.String())
		}
		cpb.Members = members
	}
	return cpb
}

func (s *ServerModel) FetchAllChats(ctx context.Context, search *chat.Search) ([]*chat.Chat, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.FetchAllChats))
	defer span.End()

	chats, err := s.Repo.ListChats(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}

	return chats, nil
}

func (s *ServerModel) FetchChatsByConversationId(ctx context.Context, id interface{}) ([]*chat.Chat, error) {

	return nil, nil
}

func (s *ServerModel) GetConversationById(context.Context, interface{}) (*conversation.Conversation, error) {

	return nil, nil
}

func (s *ServerModel) GetConversationWithMemberId(ctx context.Context, id interface{}, memberId interface{}) (*conversation.Conversation, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetConversationWithMemberId))
	defer span.End()

	cid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	mid, err := lib.ToUUID(memberId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	cvs, err := s.Repo.SelectConversation(ctx, &conversation.Search{
		Conversation: conversation.Conversation{
			BaseModel: database.BaseModel{
				ID: cid,
			},
		},
		MemberId: mid,
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	return cvs, nil
}

func (s *ServerModel) CreateChat(ctx context.Context, id interface{}, userId interface{}, payload *string, shouldSeen *bool) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CreateChat))
	defer span.End()

	userUId, err := lib.ToUUID(userId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}

	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}

	_, err = s.Repo.InsertChat(ctx, &chat.Chat{
		SenderId:       userUId,
		ConversationId: uid,
		Payload:        payload,
		Seen:           shouldSeen,
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}

	return nil
}

func (s *ServerModel) ListConversations(ctx context.Context, search *conversation.Search) ([]*conversation.Conversation, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListConversations))
	defer span.End()

	cons, err := s.Repo.ListConversations(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}

	return cons, nil
}

// Implement GetPhonePoolById
func (s *ServerModel) GetPhonePoolById(ctx context.Context, id interface{}, status *int32) (*phone_pool.PhonePool, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetPhonePoolById))
	defer span.End()

	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}

	pp, err := s.Repo.SelectPhonePool(ctx, &phone_pool.Search{
		PhonePool: phone_pool.PhonePool{
			BaseModel: database.BaseModel{
				ID: uid,
			},
			Status: status,
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}

	return pp, nil
}
