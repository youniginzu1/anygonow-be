package api

import (
	"context"

	"github.com/aqaurius6666/chatservice/src/internal/db/chat"
	"github.com/aqaurius6666/chatservice/src/internal/db/conversation"
	"github.com/aqaurius6666/chatservice/src/internal/db/notification"
	"github.com/aqaurius6666/chatservice/src/internal/lib"
	"github.com/aqaurius6666/chatservice/src/internal/lib/validate"
	"github.com/aqaurius6666/chatservice/src/internal/model"
	"github.com/aqaurius6666/chatservice/src/internal/var/c"
	"github.com/aqaurius6666/chatservice/src/internal/var/e"
	"github.com/aqaurius6666/chatservice/src/pb/chatpb"
	"github.com/aqaurius6666/chatservice/src/services/twilloclient"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/aqaurius6666/go-utils/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

type HTTPService struct {
	Model  model.Server
	Logger *logrus.Logger
}

func (s *HTTPService) FetchChats(ctx context.Context, req *chatpb.FetchPostRequest) (*chatpb.FetchOnePostResponse_Data, error) {

	return nil, nil
}

func (s *HTTPService) FetchChatById(ctx context.Context, req *chatpb.FetchOnePostRequest) (*chatpb.FetchOnePostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.FetchChatById))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId", "Id", "Timestamp", "Min"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err)
		return nil, err
	}

	cvs, err := s.Model.GetConversationWithMemberId(ctx, req.Id, req.XUserId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	chats, err := s.Model.FetchAllChats(ctx, &chat.Search{
		Chat: chat.Chat{
			ConversationId: cvs.ID,
		},
		Timestamp: &req.Timestamp,
		Min:       &req.Min,
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	if len(chats) > 0 {
		if *chats[len(chats)-1].Seen {
			return &chatpb.FetchOnePostResponse_Data{
				Chats: s.Model.ConvertChatToProtos(chats),
			}, nil
		}

		if err = s.Model.SeenChat(ctx, cvs.ID, uuid.MustParse(req.XUserId)); err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err)
			return nil, err
		}

	}
	return &chatpb.FetchOnePostResponse_Data{
		Chats: s.Model.ConvertChatToProtos(chats),
	}, nil
}

func (s *HTTPService) SendChat(ctx context.Context, req *chatpb.SendPostRequest) (*chatpb.SendPostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.SendChat))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId", "Id", "Payload"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err)
		return nil, err
	}

	cvs, err := s.Model.GetConversationWithMemberId(ctx, req.Id, req.XUserId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}

	err = s.Model.CreateChat(ctx, req.Id, req.XUserId, &req.Payload, utils.SafeBoolPtr(false))
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}

	mid, err := lib.ToUUID(req.XUserId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}

	for _, member := range cvs.Members {
		if member != mid {
			s.Model.UpsertNotification(ctx, &notification.Notification{
				UserId: member,
				Seen:   utils.BoolPtr(false),
			})
			if err != nil {
				err = xerrors.Errorf("%w", err)
				lib.RecordError(span, err)
				return nil, err
			}
		}
	}

	return &chatpb.SendPostResponse_Data{}, nil
}

func (s *HTTPService) GetConversation(ctx context.Context, req *chatpb.ConversationPostRequest) (*chatpb.ConversationPostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetConversation))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err)
		return nil, err
	}

	cIds := make([]uuid.UUID, 0)
	for _, s := range req.ConversationIds {
		cid, err := lib.ToUUID(s)
		if err != nil {
			return nil, xerrors.Errorf("%w", err)
		}
		cIds = append(cIds, cid)
	}

	cons, err := s.Model.ListConversations(ctx, &conversation.Search{
		DefaultSearchModel: database.DefaultSearchModel{
			Fields: []string{
				`"conversations"."id"`,
				`"last_chat"`,
				`"members"`,
			},
		},
		ConversationIds: cIds,
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}

	return &chatpb.ConversationPostResponse_Data{
		Conversation: s.Model.ConvertConversationToProtos(cons),
	}, nil
}

func (s *HTTPService) GetNotification(ctx context.Context, req *chatpb.NotificationGetRequest) (*chatpb.NotificationGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetNotification))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err)
		return nil, err
	}

	noti, err := s.Model.GetNotification(ctx, req.XUserId)
	if err != nil {
		// err = xerrors.Errorf("%w", err)
		// lib.RecordError(span, err)
		// return nil, err
		return &chatpb.NotificationGetResponse_Data{
			Seen: true,
		}, nil
	}

	return &chatpb.NotificationGetResponse_Data{
		Seen: utils.BoolVal(noti.Seen),
	}, nil
}

func (s *HTTPService) SeenNotification(ctx context.Context, req *chatpb.NotificationSeenPutRequest) (*chatpb.NotificationSeenPutResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.SeenNotification))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err)
		return nil, err
	}

	uid, err := lib.ToUUID(req.XUserId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}

	err = s.Model.UpsertNotification(ctx, &notification.Notification{
		UserId: uid,
		Seen:   utils.BoolPtr(true),
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}

	return &chatpb.NotificationSeenPutResponse_Data{}, nil
}

func (s *HTTPService) ProcessMessageCallback(ctx context.Context, messageCallbackData *twilloclient.TwilloMessageCallbackData) (interface{}, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ProcessMessageCallback))
	defer span.End()

	sender := messageCallbackData.From
	proxy := messageCallbackData.To
	message := messageCallbackData.Body

	conv, err := s.Model.GetConversationThroughTwillo(ctx, proxy, sender)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	for i, phone := range conv.PhoneNumberMembers {
		if phone == sender {
			err = s.Model.CreateChat(ctx, conv.ID, conv.Members[i], &message, utils.BoolPtr(false))
			if err != nil {
				err = xerrors.Errorf("%w", err)
				lib.RecordError(span, err)
				return nil, err
			}
			continue
		}

		err = s.Model.UpsertNotification(ctx, &notification.Notification{
			UserId: conv.Members[i],
			Seen:   utils.BoolPtr(false),
		})
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err)
			return nil, err
		}
	}
	return map[string]string{
		"ok": "ok",
	}, nil
}

// func (s *HTTPService) ProcessMessageCallback(ctx context.Context, messageCallbackData *twilloclient.TwilloMessageCallbackData) (interface{}, error) {
// 	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ProcessMessageCallback))
// 	defer span.End()

// 	sender := messageCallbackData.From
// 	proxy := messageCallbackData.To
// 	message := messageCallbackData.Body

// 	conv, err := s.Model.GetConversationThroughTwillo(ctx, proxy, sender)
// 	if err != nil {
// 		err = xerrors.Errorf("%w", err)
// 		lib.RecordError(span, err)
// 		return nil, err
// 	}
// 	sendFunc := func(ctx context.Context, toUserId uuid.UUID, from, to, message string) error {
// 		ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, "Goroutine-SendMessageTwillo")
// 		defer span.End()

// 		err := s.Model.SendMessage(ctx, from, to, message)
// 		if err != nil {
// 			err = xerrors.Errorf("%w", err)
// 			lib.RecordError(span, err)
// 			s.Logger.Error(err)
// 			return err
// 		}
// 		err = s.Model.UpsertNotification(ctx, &notification.Notification{
// 			UserId: toUserId,
// 			Seen:   utils.BoolPtr(false),
// 		})
// 		if err != nil {
// 			err = xerrors.Errorf("%w", err)
// 			lib.RecordError(span, err)
// 			s.Logger.Error(err)
// 			return nil
// 		}
// 		return nil
// 	}

// 	receiverPhoneNumbers := make([]string, 0)
// 	receiverUserId := make([]uuid.UUID, 0)
// 	senderIndex := 0
// 	for i, phone := range conv.PhoneNumberMembers {
// 		if phone == sender {
// 			senderIndex = i
// 			continue
// 		}
// 		receiverPhoneNumbers = append(receiverPhoneNumbers, phone)
// 		receiverUserId = append(receiverUserId, conv.Members[i])
// 	}

// 	// Check if receiver has not-seen message in app
// 	for i, userId := range receiverUserId {
// 		active, err := s.Model.IsUserActive(ctx, userId)
// 		if err != nil {
// 			err = xerrors.Errorf("%w", err)
// 			lib.RecordError(span, err)
// 			return nil, err
// 		}
// 		if active
// 		// Fetch all not-seen message
// 		chats, err := s.Model.FetchAllChats(ctx, &chat.Search{
// 			ReceiverId: userId,
// 			Chat: chat.Chat{
// 				ConversationId: conv.ID,
// 				Seen:           utils.SafeBoolPtr(false),
// 			},
// 		})
// 		if err != nil {
// 			err = xerrors.Errorf("%w", err)
// 			lib.RecordError(span, err)
// 			return nil, err
// 		}
// 		if len(chats) == 0 {
// 			continue
// 		}

// 		if !*noti.Seen {
// 			// Fetch all not-seen message
// 			chats, err := s.Model.FetchAllChats(ctx, &chat.Search{
// 				ReceiverId: userId,
// 				Chat: chat.Chat{
// 					ConversationId: conv.ID,
// 					Seen:           utils.SafeBoolPtr(false),
// 				},
// 			})
// 			if err != nil {
// 				err = xerrors.Errorf("%w", err)
// 				lib.RecordError(span, err)
// 				return nil, err
// 			}
// 			if len(chats) == 0 {
// 				continue
// 			}

// 			err = s.Model.SeenChat(ctx, conv.ID, userId)
// 			if err != nil {
// 				err = xerrors.Errorf("%w", err)
// 				lib.RecordError(span, err)
// 				return nil, err
// 			}

// 			err = s.Model.UpsertNotification(ctx, &notification.Notification{
// 				UserId: userId,
// 				Seen:   utils.BoolPtr(true),
// 			})
// 			if err != nil {
// 				err = xerrors.Errorf("%w", err)
// 				lib.RecordError(span, err)
// 				return nil, err
// 			}

// 			combinedMessage := s.Model.CombineMessage(ctx, chats)

// 			err = sendFunc(ctx, userId, proxy, receiverPhoneNumbers[i], combinedMessage)
// 			if err != nil {
// 				err = xerrors.Errorf("%w", err)
// 				lib.RecordError(span, err)
// 				s.Logger.Error(err)
// 				return nil, err
// 			}
// 		}

// 	}

// 	err = s.Model.CreateChat(ctx, conv.ID, conv.Members[senderIndex], &message, utils.BoolPtr(false))
// 	if err != nil {
// 		err = xerrors.Errorf("%w", err)
// 		lib.RecordError(span, err)
// 		return nil, err
// 	}
// 	return map[string]string{
// 		"ok": "ok",
// 	}, nil
// }
