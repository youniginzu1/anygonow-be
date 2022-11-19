package cockroach

import (
	"context"

	"github.com/aqaurius6666/chatservice/src/internal/db/conversation"
	"github.com/aqaurius6666/chatservice/src/internal/lib"
	"github.com/aqaurius6666/chatservice/src/internal/var/c"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func applySearchConversation(db *gorm.DB, search *conversation.Search) *gorm.DB {
	if search.ID != uuid.Nil {
		db = db.Where(conversation.Conversation{
			BaseModel: database.BaseModel{
				ID: search.ID,
			},
		})
	}
	if search.MemberId != uuid.Nil {
		db = db.Where(`cast(? as uuid) = any ("conversations"."members" :: uuid[])`, search.MemberId)
	}
	if search.Members != nil {
		db = db.Where(conversation.Conversation{
			Members: search.Members,
		})
	}
	if search.ConversationIds != nil {
		db = db.Where(`"conversations"."id" = any(? :: uuid[])`, search.ConversationIds)
	}
	if search.OrderId != uuid.Nil {
		db = db.Where(conversation.Conversation{
			OrderId: search.OrderId,
		})
	}
	if search.ServiceId != uuid.Nil {
		db = db.Where(conversation.Conversation{
			ServiceId: search.ServiceId,
		})
	}
	if search.MemberPhone != nil {
		db = db.Where(`cast(? as varchar(16)) = any ("conversations"."phone_number_members" ::varchar(16)[])`, search.MemberPhone)
	}
	if search.PhonePoolId != uuid.Nil {
		db = db.Where(conversation.Conversation{
			PhonePoolId: search.PhonePoolId,
		})
	}

	return db
}

func (s *ServerCDBRepo) SelectConversation(ctx context.Context, search *conversation.Search) (*conversation.Conversation, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.SelectConversation))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r := conversation.Conversation{}
	if err := applySearchConversation(s.Db, search).WithContext(ctx).First(&r).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = conversation.ErrNotFound
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	return &r, nil
}

func (s *ServerCDBRepo) InsertConversation(ctx context.Context, value *conversation.Conversation) (*conversation.Conversation, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.InsertConversation))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := s.Db.WithContext(ctx).Create(value).Error; err != nil {
		err = xerrors.Errorf("%w", conversation.ErrInsertFail)
		lib.RecordError(span, err)
		return nil, err
	}
	return value, nil
}

func (s *ServerCDBRepo) ListConversations(ctx context.Context, search *conversation.Search) ([]*conversation.Conversation, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListConversations))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	db := s.Db
	r := make([]*conversation.Conversation, 0)
	if err := applySearchConversation(db, search).WithContext(ctx).Select(search.Fields).
		Joins(`left join (?) as chat on "chat"."conversation_id" = "conversations"."id"`, getSubQueryList(db)).
		Group(`"conversations"."id" , "last_chat"`).
		Order(`"last_chat" DESC`).
		Find(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	return r, nil
}

func getSubQueryList(db *gorm.DB) *gorm.DB {
	db = db.Table("chats")
	db = db.Select(`"chats"."conversation_id", MAX("chats"."updated_at") as last_chat`)
	db = db.Group(`"chats"."conversation_id"`)
	return db
}

func (s *ServerCDBRepo) UpdateConversation(ctx context.Context, search *conversation.Search, value *conversation.Conversation) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UpdateConversation))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := applySearchConversation(s.Db, search).WithContext(ctx).Updates(value).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	return nil
}
func (s *ServerCDBRepo) SetConversationPhonePoolNull(ctx context.Context, search *conversation.Search) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.SetConversationPhonePoolNull))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := applySearchConversation(s.Db, search).WithContext(ctx).Update("phone_pool_id = ?", uuid.Nil).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	return nil
}
func (s *ServerCDBRepo) ListUnusedConversation(ctx context.Context, search *conversation.Search) ([]*conversation.Conversation, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListUnusedConversation))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ret := make([]*conversation.Conversation, 0)
	exprs := make([]clause.Expression, 0)
	for _, p := range search.PhoneNumberMembers {
		exprs = append(exprs, clause.Expr{
			SQL:  `cast(? as varchar(16)) = any(phone_number_members)`,
			Vars: []interface{}{p},
		})
	}

	if err := s.Db.WithContext(ctx).Table("conversations").Where(clause.AndConditions{
		Exprs: []clause.Expression{
			clause.Eq{
				Column: "status",
				Value:  1,
			},
			clause.AndConditions{
				Exprs: exprs,
			},
			clause.Not(
				clause.Eq{
					Column: "phone_pool_id",
					Value:  uuid.Nil,
				},
			),
		},
	}).Debug().Find(&ret).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	return ret, nil
}
