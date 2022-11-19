package model

import (
	"context"

	"github.com/aqaurius6666/apiservice/src/internal/db/business"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/apiservice/src/pb"
	"github.com/aqaurius6666/go-utils/utils"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

var (
	_ ChatModel = (*ServerModel)(nil)
)

type ChatModel interface {
	NewConversation(context.Context, uuid.UUID, uuid.UUID, uuid.UUID) (uuid.UUID, error)
	GetConversation(context.Context, string, []string) (*pb.ConversationPostResponse_Data, error)
}

func (s *ServerModel) NewConversation(ctx context.Context, bid, cid, oid uuid.UUID) (uuid.UUID, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.NewConversation))
	defer span.End()

	bus, err := s.GetBusinessById(ctx, bid)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return uuid.Nil, err
	}

	order, err := s.GetOrderById(ctx, oid)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return uuid.Nil, err
	}

	members := []string{cid.String(), bid.String()}
	phones := []string{*order.CustomerPhone, *bus.Phone}

	convId, err := s.Chat.NewConversation(ctx, oid.String(), order.ServiceId.String(), members, phones)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return uuid.Nil, err
	}
	id, err := lib.ToUUID(convId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return uuid.Nil, err
	}
	return id, nil
}

func (s *ServerModel) GetConversation(ctx context.Context, userId string, convIds []string) (*pb.ConversationPostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetConversation))
	defer span.End()

	res, err := s.Chat.GetConversation(ctx, userId, convIds)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	memberIds := make([]string, 0)
	for _, c := range res.Conversation {
		memberIds = append(memberIds, c.Members...)
	}
	idNames, err := s.Repo.GetMapIdName(ctx, &business.Search{
		BothId: memberIds,
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	convs := make([]*pb.Conversation, 0)
	for _, r := range res.Conversation {
		mem := make([]*pb.Conversation_Member, 0)
		for _, m := range r.Members {
			mem = append(mem, &pb.Conversation_Member{
				Id:   m,
				Name: utils.StrVal(idNames[m]),
			})
		}
		convs = append(convs, &pb.Conversation{
			Members:  mem,
			Id:       r.Id,
			LastChat: r.LastChat,
		})
	}
	return &pb.ConversationPostResponse_Data{
		Conversation: convs,
	}, nil
}
