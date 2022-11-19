package api

import (
	"context"

	"github.com/aqaurius6666/chatservice/src/internal/db/chat"
	"github.com/aqaurius6666/chatservice/src/internal/db/conversation"
	"github.com/aqaurius6666/chatservice/src/internal/lib"
	"github.com/aqaurius6666/chatservice/src/internal/var/c"
	"github.com/aqaurius6666/chatservice/src/pb/chatpb"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/aqaurius6666/go-utils/utils"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

func (s *ApiServer) Chat(stream chatpb.ChatService_ChatServer) error {
	panic("not implemented")
	// fmt.Println("Start new")
	// s.handleStream(stream)
	// return nil
}

func (s *ApiServer) GetConversation(ctx context.Context, req *chatpb.ConversationPostRequest) (*chatpb.ConversationPostResponse, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetConversation))
	defer span.End()
	res, err := s.HTTPService.GetConversation(ctx, req)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	return &chatpb.ConversationPostResponse{
		Data: res,
	}, nil
}

func (s *ApiServer) NewConversation(ctx context.Context, req *chatpb.NewConversationRequest) (*chatpb.NewConversationResponse, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.NewConversation))
	defer span.End()

	oid, err := lib.ToUUID(req.OrderId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		panic(err)
	}
	sid, err := lib.ToUUID(req.ServiceId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		panic(err)
	}

	uuids := make([]uuid.UUID, 0)
	for _, id := range req.MemberIds {
		i, err := lib.ToUUID(id)
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err)
			panic(err)
		}
		uuids = append(uuids, i)
	}

	sortedUuid, sortedPhone := lib.SortUUIDWithPhoneNumber(uuids, req.PhoneNumbers)

	conv, err := s.Model.GetConversationFullMember(ctx, oid, sid, sortedUuid)

	if err == nil {
		// conversation existed
		if utils.Int32Val(conv.Status) == 0 {
			// Open conversation
			return &chatpb.NewConversationResponse{
				Id: conv.ID.String(),
			}, nil
		}

		ok, err := s.Model.ValidateConversationPhonePool(ctx, conv.PhonePoolId, sortedPhone)
		if err == nil && ok {
			if err = s.Repo.UpdateConversation(ctx, &conversation.Search{
				Conversation: conversation.Conversation{BaseModel: database.BaseModel{ID: conv.ID}},
			}, &conversation.Conversation{
				Status: utils.Int32Ptr(0),
			}); err != nil {
				err = xerrors.Errorf("%w", err)
				lib.RecordError(span, err)
				panic(err)
			}
			return &chatpb.NewConversationResponse{
				Id: conv.ID.String(),
			}, nil
		}
		ppId, err := s.Model.GetReleaseBuyPhonePool(ctx, sortedPhone)
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err)
			panic(err)
		}
		if err = s.Repo.UpdateConversation(ctx, &conversation.Search{
			Conversation: conversation.Conversation{BaseModel: database.BaseModel{ID: conv.ID}},
		}, &conversation.Conversation{
			Status:      utils.Int32Ptr(0),
			PhonePoolId: ppId,
		}); err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err)
			panic(err)
		}
		return &chatpb.NewConversationResponse{
			Id: conv.ID.String(),
		}, nil
	}

	// conversation not found
	ppId, err := s.Model.GetReleaseBuyPhonePool(ctx, sortedPhone)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		panic(err)
	}
	cvs, err := s.Model.NewConversation(ctx, oid, sid, sortedUuid, sortedPhone, ppId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		panic(err)
	}
	return &chatpb.NewConversationResponse{
		Id: cvs.ID.String(),
	}, nil

	// conversation existed with status = 1

	// update conversation

	// handleWhenPhonePoolAvailableFunc := func(ctx context.Context, pp *phone_pool.PhonePool) (*chatpb.NewConversationResponse, error) {
	// 	cvs, err := s.Model.GetConversationFullMember(ctx, oid, sid, sortedUuid)
	// 	if err == nil {
	// 		// exist conversation
	// 		if utils.Int32Val(cvs.Status) == 0 {
	// 			// Open conversation
	// 			err = xerrors.Errorf("%w", e.ErrOpenConversation)
	// 			lib.RecordError(span, err)
	// 			panic(err)
	// 		}
	// 		var phoneUpdate []string
	// 		var phonePoolUpdate uuid.UUID

	// 		if !lib.IsSliceEqual(sortedPhone, cvs.PhoneNumberMembers) {
	// 			// phone number change
	// 			phoneUpdate = sortedPhone
	// 		}

	// 		if cvs.PhonePoolId == uuid.Nil {
	// 			phonePoolUpdate = pp.ID
	// 		}

	// 		err = s.Model.UpdatePhoneConversation(ctx, cvs.ID, phoneUpdate, phonePoolUpdate)
	// 		if err != nil {
	// 			err = xerrors.Errorf("%w", err)
	// 			lib.RecordError(span, err)
	// 			return nil, err
	// 		}

	// 		return &chatpb.NewConversationResponse{
	// 			Id: cvs.ID.String(),
	// 		}, nil

	// 	}
	// 	cvs, err = s.Model.NewConversation(ctx, oid, sid, sortedUuid, sortedPhone, pp.ID)
	// 	if err != nil {
	// 		err = xerrors.Errorf("%w", err)
	// 		lib.RecordError(span, err)
	// 		panic(err)
	// 	}
	// 	return &chatpb.NewConversationResponse{
	// 		Id: cvs.ID.String(),
	// 	}, nil
	// }
	// phonePool, err := s.Model.GetAvailablePhone(ctx, sortedPhone)
	// if err == nil {
	// 	return handleWhenPhonePoolAvailableFunc(ctx, phonePool)
	// }
	// if err != nil && !errors.Is(err, e.ErrNoPhoneAvailable) {
	// 	err = xerrors.Errorf("%w", err)
	// 	lib.RecordError(span, err)
	// 	return nil, err
	// }
	// err = s.Model.SyncPhonePool(ctx)

	// if err != nil {
	// 	err = xerrors.Errorf("%w", err)
	// 	lib.RecordError(span, err)
	// 	return nil, err
	// }
	// phonePool, err = s.Model.GetAvailablePhone(ctx, sortedPhone)
	// if err == nil {
	// 	return handleWhenPhonePoolAvailableFunc(ctx, phonePool)
	// }
	// if err != nil && !errors.Is(err, e.ErrNoPhoneAvailable) {
	// 	err = xerrors.Errorf("%w", err)
	// 	lib.RecordError(span, err)
	// 	return nil, err
	// }
	// phonePool, err = s.Model.BuyNewPhone(ctx)
	// if err != nil {
	// 	err = xerrors.Errorf("%w", err)
	// 	lib.RecordError(span, err)
	// 	return nil, err
	// }
	// return handleWhenPhonePoolAvailableFunc(ctx, phonePool)
}

func (s *ApiServer) TriggerSendSMS(ctx context.Context, req *chatpb.TriggerSendSMSRequest) (*chatpb.TriggerSendSMSResponse, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.TriggerSendSMS))
	defer span.End()
	userId := uuid.MustParse(req.UserId)

	convs, err := s.Model.ListConversations(ctx, &conversation.Search{
		DefaultSearchModel: database.DefaultSearchModel{
			Fields: []string{
				`"conversations"."id"`,
				`"last_chat"`,
				`"members"`,
				`"phone_number_members"`,
				`"phone_pool_id"`,
			},
		},
		MemberId: userId,
		Conversation: conversation.Conversation{
			Status: utils.Int32Ptr(0),
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}

	conversationHandle := func(ctx context.Context, conv *conversation.Conversation, userId uuid.UUID) error {
		ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, "Goroutine-conversationHandle")
		defer span.End()

		chats, err := s.Model.FetchAllChats(ctx, &chat.Search{
			ReceiverId: userId,
			Chat: chat.Chat{
				ConversationId: conv.ID,
				Seen:           utils.SafeBoolPtr(false),
			},
		})
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err)
			return err
		}
		if len(chats) == 0 {
			return nil
		}
		phonePool, err := s.Model.GetPhonePoolById(ctx, conv.PhonePoolId, utils.Int32Ptr(0))
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err)
			return err
		}
		combined := s.Model.CombineMessage(ctx, chats)

		index := 0
		for i := 0; i < len(conv.Members); i++ {
			if conv.Members[i] == userId {
				index = i
				break
			}
		}
		err = s.Model.SendMessage(ctx, lib.StandardPhoneNumber(*phonePool.PhoneNumber), conv.PhoneNumberMembers[index], combined)
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err)
			return err
		}
		err = s.Model.SeenChat(ctx, conv.ID, userId)
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err)
			return err
		}
		return nil
	}

	for _, conv := range convs {
		go func(conv *conversation.Conversation) {
			if err := conversationHandle(context.TODO(), conv, userId); err != nil {
				s.Logger.Error(err)
			}
		}(conv)
	}
	err = s.Model.SetUserInactiveTimeout(ctx, userId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}

	return &chatpb.TriggerSendSMSResponse{}, nil
}

func (s *ApiServer) CloseConversation(ctx context.Context, req *chatpb.CloseConversationRequest) (*chatpb.CloseConversationResponse, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CloseConversation))
	defer span.End()

	orderId := uuid.MustParse(req.OrderId)

	err := s.Model.CloseConversation(ctx, orderId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		panic(err)
	}
	return &chatpb.CloseConversationResponse{}, nil
}
