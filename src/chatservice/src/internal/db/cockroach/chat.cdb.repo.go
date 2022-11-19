package cockroach

import (
	"context"

	"github.com/aqaurius6666/chatservice/src/internal/db/chat"
	"github.com/aqaurius6666/chatservice/src/internal/lib"
	"github.com/aqaurius6666/chatservice/src/internal/var/c"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func applySearchChat(db *gorm.DB, search *chat.Search) *gorm.DB {
	if search.ID != uuid.Nil {
		db = db.Where(chat.Chat{
			BaseModel: database.BaseModel{
				ID: search.ID,
			},
		})
	}

	if search.Timestamp != nil {
		db = db.Where(`"chats"."created_at" < ?`, *search.Timestamp)
	}
	if search.Min != nil {
		db = db.Limit(int(*search.Min))
	}
	if search.ConversationId != uuid.Nil {
		db = db.Where(chat.Chat{
			ConversationId: search.ConversationId,
		})
	}
	if search.ReceiverId != uuid.Nil {
		db = db.Where(clause.Neq{
			Column: "sender_id",
			Value:  search.ReceiverId,
		})
	}
	if search.Seen != nil {
		db = db.Where(chat.Chat{
			Seen: search.Seen,
		})
	}
	db = db.Order("created_at ASC")

	return db
}

func (s *ServerCDBRepo) SelectChat(ctx context.Context, search *chat.Search) (*chat.Chat, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.SelectChat))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r := chat.Chat{}
	if err := applySearchChat(s.Db, search).WithContext(ctx).First(&r).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = chat.ErrNotFound
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	return &r, nil
}

func (s *ServerCDBRepo) InsertChat(ctx context.Context, value *chat.Chat) (*chat.Chat, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.InsertChat))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := s.Db.WithContext(ctx).Create(value).Error; err != nil {
		err = xerrors.Errorf("%w", chat.ErrInsertFail)
		lib.RecordError(span, err)
		return nil, err
	}
	return value, nil
}

func (s *ServerCDBRepo) ListChats(ctx context.Context, search *chat.Search) ([]*chat.Chat, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListChats))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	r := make([]*chat.Chat, 0)
	if err := applySearchChat(s.Db, search).WithContext(ctx).Find(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	return r, nil
}

func (s *ServerCDBRepo) UpdateChat(ctx context.Context, search *chat.Search, value *chat.Chat) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UpdateChat))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := applySearchChat(s.Db, search).WithContext(ctx).Updates(value).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	return nil
}
