package model

import (
	"bytes"
	"context"
	"html/template"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aqaurius6666/chatservice/src/internal/db/chat"
	"github.com/aqaurius6666/chatservice/src/internal/db/conversation"
	"github.com/aqaurius6666/chatservice/src/internal/db/phone_pool"
	"github.com/aqaurius6666/chatservice/src/internal/lib"
	"github.com/aqaurius6666/chatservice/src/internal/lib/unleash"
	"github.com/aqaurius6666/chatservice/src/internal/var/c"
	"github.com/aqaurius6666/chatservice/src/internal/var/e"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/aqaurius6666/go-utils/utils"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

type Conversation interface {
	NewConversation(ctx context.Context, orderId uuid.UUID, serviceId uuid.UUID, memberIds []uuid.UUID, phones []string, phonePoolId uuid.UUID) (*conversation.Conversation, error)
	GetConversationFullMember(ctx context.Context, orderId uuid.UUID, serviceId uuid.UUID, memberIds []uuid.UUID) (*conversation.Conversation, error)
	GetConversationThroughTwillo(ctx context.Context, proxyAddress string, member string) (*conversation.Conversation, error)
	UpdatePhoneConversation(ctx context.Context, id uuid.UUID, phones []string, phonePoolId uuid.UUID, orderId uuid.UUID) error

	SendMessage(ctx context.Context, from string, to string, message string) error
	GetOrBuyPhonePool(ctx context.Context, serviceId uuid.UUID, phones []string) (*phone_pool.PhonePool, error)
	CombineMessage(ctx context.Context, chats []*chat.Chat) string
	ReleaseConversationPhonePool(ctx context.Context, phones []string) (uuid.UUID, error)

	BuyNewPhone(ctx context.Context) (*phone_pool.PhonePool, error)
	ListResourcePhone(ctx context.Context) (phones []string, sids []string, err error)
	CloseConversation(ctx context.Context, orderId uuid.UUID) error

	UpsertConversation(ctx context.Context, id uuid.UUID, orderId uuid.UUID, serviceId uuid.UUID, memberIds []uuid.UUID, phones []string, phonePoolId uuid.UUID) (*conversation.Conversation, error)
	SyncPhonePool(ctx context.Context) error
	GetAvailablePhone(ctx context.Context, phones []string) (*phone_pool.PhonePool, error)
	GetAvailablePhoneWithSync(ctx context.Context, phones []string) (*phone_pool.PhonePool, error)
	ShouldReuseConversationPhonePool(ctx context.Context, conversationID uuid.UUID, phones []string, phonePoolId uuid.UUID) (bool, error)
	GetReleaseBuyPhonePool(ctx context.Context, phones []string) (uuid.UUID, error)
	ValidateConversationPhonePool(ctx context.Context, phonePoolId uuid.UUID, phones []string) (bool, error)
}

type phonePool struct {
	Id          uuid.UUID
	PhoneNumber string
	sid         string
	status      int32
}

type sortablePhone []phonePool

func (a sortablePhone) Len() int           { return len(a) }
func (a sortablePhone) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortablePhone) Less(i, j int) bool { return a[i].PhoneNumber < a[j].PhoneNumber }
func (a sortablePhone) Search(phone string) int {
	if len(a) == 0 {
		return -1
	}
	return sort.Search(len(a), func(i int) bool {
		return a[i].PhoneNumber >= phone
	})
}

func (s *ServerModel) GetReleaseBuyPhonePool(ctx context.Context, sortedPhone []string) (uuid.UUID, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetReleaseBuyPhonePool))
	defer span.End()

	pp, err := s.GetAvailablePhoneWithSync(ctx, sortedPhone)
	if err == nil {
		return pp.ID, nil
	}
	pps, err := s.ListPhonePool(ctx, &phone_pool.Search{
		PhonePool: phone_pool.PhonePool{Status: utils.Int32Ptr(0)},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return uuid.Nil, err
	}
	maxPhoneNumber, _ := strconv.Atoi(unleash.GetVariant("chatservice.twillo.max-phone-number"))
	var pid uuid.UUID
	if len(pps) >= maxPhoneNumber {
		pid, err = s.ReleaseConversationPhonePool(ctx, sortedPhone)
		if err != nil {
			err = xerrors.Errorf("%w", e.ErrMaxPhoneNumbers)
			lib.RecordError(span, err)
			return uuid.Nil, err
		}
	} else {
		phonePool, err := s.BuyNewPhone(ctx)
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err)
			return uuid.Nil, err
		}
		pid = phonePool.ID
	}
	return pid, nil
}

func (s *ServerModel) UpsertConversation(ctx context.Context, id uuid.UUID, orderId uuid.UUID, serviceId uuid.UUID, memberIds []uuid.UUID, phones []string, phonePoolId uuid.UUID) (*conversation.Conversation, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UpsertConversation))
	defer span.End()
	if id == uuid.Nil {
		conv, err := s.NewConversation(ctx, orderId, serviceId, memberIds, phones, phonePoolId)
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err)
			return nil, err
		}
		return conv, nil
	}
	err := s.UpdatePhoneConversation(ctx, id, phones, phonePoolId, orderId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	return &conversation.Conversation{
		BaseModel: database.BaseModel{
			ID: id,
		},
		OrderId:            orderId,
		ServiceId:          serviceId,
		Members:            memberIds,
		PhoneNumberMembers: phones,
		PhonePoolId:        phonePoolId,
	}, nil
}

func (s *ServerModel) ShouldReuseConversationPhonePool(ctx context.Context, convId uuid.UUID, phones []string, phonePoolId uuid.UUID) (bool, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ShouldReuseConversationPhonePool))
	defer span.End()
	pp, err := s.Repo.ListConversations(ctx, &conversation.Search{
		Conversation: conversation.Conversation{
			Status:      utils.Int32Ptr(0),
			PhonePoolId: phonePoolId,
		},
	})
	if err == nil || len(pp) > 0 {
		return false, nil
	}
	pp, err = s.Repo.ListUnusedConversation(ctx, &conversation.Search{
		Conversation: conversation.Conversation{
			PhoneNumberMembers: phones,
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return false, err
	}

	if len(pp) == 1 || len(pp) == 0 {
		return true, nil
	}
	return false, nil
}

func (s *ServerModel) GetAvailablePhoneWithSync(ctx context.Context, phones []string) (*phone_pool.PhonePool, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetAvailablePhoneWithSync))
	defer span.End()
	pp, err := s.GetAvailablePhone(ctx, phones)
	if err == nil {
		return pp, nil
	}
	err = s.SyncPhonePool(ctx)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	pp, err = s.GetAvailablePhone(ctx, phones)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	return pp, nil
}

// Implement ReleaseConversationPhonePool
func (s *ServerModel) ReleaseConversationPhonePool(ctx context.Context, phones []string) (uuid.UUID, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ReleaseConversationPhonePool))
	defer span.End()

	convs, err := s.Repo.ListUnusedConversation(ctx, &conversation.Search{
		Conversation: conversation.Conversation{
			PhoneNumberMembers: phones,
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return uuid.Nil, err
	}
	if len(convs) == 0 {
		err = xerrors.Errorf("%w", e.ErrNoConversationAvailable)
		lib.RecordError(span, err)
		return uuid.Nil, err
	}
	conv := convs[rand.Intn(len(convs))]
	err = s.Repo.SetConversationPhonePoolNull(ctx, &conversation.Search{
		Conversation: conversation.Conversation{BaseModel: database.BaseModel{ID: conv.ID}},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return uuid.Nil, err
	}
	return conv.PhonePoolId, nil
}

// Implement GetOrBuyPhonePool
func (s *ServerModel) GetOrBuyPhonePool(ctx context.Context, serviceId uuid.UUID, phones []string) (*phone_pool.PhonePool, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetOrBuyPhonePool))
	defer span.End()
	err := s.SyncPhonePool(ctx)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	pp, err := s.GetAvailablePhone(ctx, phones)
	if err == nil {
		return pp, nil
	}

	pp, err = s.BuyNewPhone(ctx)
	if err != nil {
		err = xerrors.Errorf("%w", e.ErrBuyPhoneFail)
		lib.RecordError(span, err)
		return nil, err
	}

	return pp, nil
}

// Implement CloseConversation
func (s *ServerModel) CloseConversation(ctx context.Context, orderId uuid.UUID) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CloseConversation))
	defer span.End()
	err := s.Repo.UpdateConversation(ctx, &conversation.Search{
		Conversation: conversation.Conversation{
			OrderId: orderId,
			Status:  utils.Int32Ptr(0),
		},
	}, &conversation.Conversation{
		Status: utils.Int32Ptr(1),
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	return nil
}

// Implement SyncPhonePool
func (s *ServerModel) SyncPhonePool(ctx context.Context) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.SyncPhonePool))
	defer span.End()
	// defer func() {
	// 	if ierr := recover(); ierr != nil {
	// 		s.Logger.Error(ierr)
	// 		err = ierr.(error)
	// 	}
	// }()
	// sync phone pool from twillo to db
	twilloPhonePool, twilloPhoneSids, err := s.ListResourcePhone(ctx)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}

	dbPhonePool, err := s.ListPhonePool(ctx, &phone_pool.Search{
		DefaultSearchModel: database.DefaultSearchModel{
			OrderBy:   "phone_number",
			OrderType: "asc",
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	myPP := make(sortablePhone, 0)
	for _, pp := range dbPhonePool {
		myPP = append(myPP, phonePool{
			PhoneNumber: *pp.PhoneNumber,
			sid:         *pp.Sid,
			status:      *pp.Status,
			Id:          pp.ID,
		})
	}
	if !sort.IsSorted(myPP) {
		if err != nil {
			err = xerrors.New("phone pool is not sorted")
			lib.RecordError(span, err)
			return err
		}
	}

	mapPhonePoolSync := make(map[string]bool)

	for i, phone := range twilloPhonePool {
		findIndex := myPP.Search(phone)
		if findIndex == -1 {
			// Not found, insert into phone pool
			_, err := s.Repo.InsertPhonePool(ctx, &phone_pool.PhonePool{
				PhoneNumber: utils.StrPtr(phone),
				Sid:         utils.StrPtr(twilloPhoneSids[i]),
				Status:      utils.Int32Ptr(0),
			})
			if err != nil {
				err = xerrors.Errorf("%w", err)
				lib.RecordError(span, err)
				return err
			}
		} else {
			// Found, check status whether equal 0
			if myPP[findIndex].status != 0 {
				// Not equal 0, update status to 0
				err := s.Repo.UpdatePhonePool(ctx, &phone_pool.Search{
					PhonePool: phone_pool.PhonePool{
						BaseModel: database.BaseModel{
							ID: myPP[findIndex].Id,
						},
					},
				}, &phone_pool.PhonePool{
					Status: utils.Int32Ptr(0),
				})
				if err != nil {
					err = xerrors.Errorf("%w", err)
					lib.RecordError(span, err)
					return err
				}
			}
		}
		mapPhonePoolSync[phone] = true
	}
	notSyncPhone := make([]phonePool, 0)
	for _, phone := range myPP {
		if _, ok := mapPhonePoolSync[phone.PhoneNumber]; !ok {
			// Not found, update status to 1 in phone pool
			notSyncPhone = append(notSyncPhone, phone)

			if err := s.Repo.UpdatePhonePool(ctx, &phone_pool.Search{
				PhonePool: phone_pool.PhonePool{
					BaseModel: database.BaseModel{
						ID: phone.Id,
					},
				},
			}, &phone_pool.PhonePool{
				Status: utils.Int32Ptr(1),
			}); err != nil {
				err = xerrors.Errorf("%w", err)
				lib.RecordError(span, err)
				return err
			}

		}
	}
	if len(notSyncPhone) > 0 {
		func() {
			if err = s.Mail.SendMail(c.ADMIN_MAIL, s.buildMailNotSyncedTemplate(c.ADMIN_MAIL, notSyncPhone)); err != nil {
				panic(err)
			}
		}()
	}

	return nil
}

// Implement ListResourcePhone
func (s *ServerModel) ListResourcePhone(ctx context.Context) ([]string, []string, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListResourcePhone))
	defer span.End()
	phones, sids, err := s.Twillo.ListResourcePhone(ctx)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, nil, err
	}
	return phones, sids, nil
}

func (s *ServerModel) BuyNewPhone(ctx context.Context) (*phone_pool.PhonePool, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.BuyNewPhone))
	defer span.End()
	var phone *string
	var sid *string
	var err error
	phone, err = s.Twillo.ListAvailablePhoneNumber(ctx)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}

	phone, sid, err = s.Twillo.BuyPhoneNumber(ctx, phone)
	if err != nil {
		err = xerrors.Errorf("%w", e.ErrBuyPhoneFail)
		lib.RecordError(span, err)
		return nil, err
	}
	phonePool, err := s.Repo.InsertPhonePool(ctx, &phone_pool.PhonePool{
		PhoneNumber: phone,
		Sid:         sid,
		Status:      utils.Int32Ptr(0),
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	go func(phonePool phone_pool.PhonePool) {
		s.Mail.SendMail(c.ADMIN_MAIL, s.buildMailBuyPhoneTemplate(c.ADMIN_MAIL, phonePool))
	}(*phonePool)
	return phonePool, nil
}

func (s *ServerModel) CombineMessage(ctx context.Context, chats []*chat.Chat) string {
	_, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CombineMessage))
	defer span.End()

	payloads := make([]string, len(chats))
	for _, chat := range chats {
		payloads = append(payloads, *chat.Payload)
	}

	return strings.Join(payloads, "\n")
}
func (s *ServerModel) ValidateConversationPhonePool(ctx context.Context, phonePoolId uuid.UUID, phones []string) (bool, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ValidateConversationPhonePool))
	defer span.End()

	pool, err := s.Repo.ListAvailablePhone(ctx, phones)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return false, err
	}
	for _, phone := range pool {
		if phone.ID == phonePoolId {
			return true, nil
		}
	}

	return false, nil
}
func (s *ServerModel) GetAvailablePhone(ctx context.Context, phones []string) (*phone_pool.PhonePool, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetAvailablePhone))
	defer span.End()

	pool, err := s.Repo.ListAvailablePhone(ctx, phones)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	if len(pool) == 0 {
		err = xerrors.Errorf("%w", e.ErrNoPhoneAvailable)
		lib.RecordError(span, err)
		return nil, err
	}

	return pool[rand.Int()%len(pool)], nil
}

func (s *ServerModel) GetConversationThroughTwillo(ctx context.Context, proxyAddress string, member string) (*conversation.Conversation, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetConversationThroughTwillo))
	defer span.End()

	phonePool, err := s.Repo.SelectPhonePool(ctx, &phone_pool.Search{
		PhonePool: phone_pool.PhonePool{
			PhoneNumber: &proxyAddress,
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	conv, err := s.Repo.SelectConversation(ctx, &conversation.Search{
		MemberPhone: &member,
		Conversation: conversation.Conversation{
			PhonePoolId: phonePool.ID,
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	return conv, nil
}

func (s *ServerModel) SendMessage(ctx context.Context, from, to, message string) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.SendMessage))
	defer span.End()

	err := s.Twillo.SendMessage(ctx, from, to, message)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	return nil
}

func (s *ServerModel) NewConversation(ctx context.Context, orderId uuid.UUID, serviceId uuid.UUID, memberIds []uuid.UUID, phones []string, phonePoolId uuid.UUID) (*conversation.Conversation, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.NewConversation))
	defer span.End()

	cvs, err := s.Repo.InsertConversation(ctx, &conversation.Conversation{
		OrderId:            orderId,
		Members:            memberIds,
		PhoneNumberMembers: phones,
		PhonePoolId:        phonePoolId,
		ServiceId:          serviceId,
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	return cvs, nil
}

func (s *ServerModel) GetConversationFullMember(ctx context.Context, orderId uuid.UUID, serviceId uuid.UUID, memberIds []uuid.UUID) (*conversation.Conversation, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetConversationFullMember))
	defer span.End()

	cvs, err := s.Repo.SelectConversation(ctx, &conversation.Search{
		Conversation: conversation.Conversation{
			// OrderId:   orderId,
			ServiceId: serviceId,
			Members:   memberIds,
			// Status:    utils.Int32Ptr(0),
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	return cvs, err
}
func (s *ServerModel) UpdatePhoneConversation(ctx context.Context, id uuid.UUID, phones []string, phonePoolId uuid.UUID, orderId uuid.UUID) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UpdatePhoneConversation))
	defer span.End()

	err := s.Repo.UpdateConversation(ctx, &conversation.Search{
		Conversation: conversation.Conversation{
			BaseModel: database.BaseModel{
				ID: id,
			},
		},
	}, &conversation.Conversation{
		OrderId:            orderId,
		PhoneNumberMembers: phones,
		PhonePoolId:        phonePoolId,
		Status:             utils.Int32Ptr(0),
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	return nil
}

func (s *ServerModel) buildMailNotSyncedTemplate(to string, phones []phonePool) []byte {

	tmpl, err := template.New("mail_template").Parse(`Subject: [AnyGoNow] Sync phone number
From: Anygonow Service <no-reply@anygonow.com>
To: {{ .To }}
MIME-Version: 1.0
Content-Type: text/html; charset=UTF-8

<html>
<body>
<p>
	<b>List phone numbers have not been synced yet:</b>
</p>
<p>
	{{ range $i, $phone := .Phones}}
		<b>{{ $phone.PhoneNumber }}</b> - <b>{{ $phone.Id }}</b>
		<br/>
	{{ end }}
</p>
<p>
	<b>Time sync:</b> ` + time.Now().Format(time.RFC1123) + `
</p>

</body>
</html>																											
`)
	if err != nil {
		panic(err)
	}
	out := &bytes.Buffer{}
	err = tmpl.Execute(out, struct {
		To     string
		Phones []phonePool
	}{
		To:     to,
		Phones: phones,
	})
	if err != nil {
		panic(err)
	}
	return out.Bytes()
}

func (s *ServerModel) buildMailBuyPhoneTemplate(to string, phone phone_pool.PhonePool) []byte {

	tmpl, err := template.New("mail_template").Parse(`Subject: [AnyGoNow] Buy new phone
From: Anygonow Service <no-reply@anygonow.com>
To: {{ .To }}
MIME-Version: 1.0
Content-Type: text/html; charset=UTF-8

<html>
<body>
<p>
	<b>System automatically bought new phone.</b>
</p>
<p>
	<b>Phone number:</b> {{ .PhoneNumber }}<br/>
	<b>Phone pool id:</b> {{ .PhonePoolId }}
</p>
<p>
	<b>Time sync:</b> ` + time.Now().Format(time.RFC1123) + `
</p>

</body>
</html>																											
`)
	if err != nil {
		panic(err)
	}
	out := &bytes.Buffer{}
	err = tmpl.Execute(out, struct {
		To          string
		PhoneNumber string
		PhonePoolId string
	}{
		To:          to,
		PhoneNumber: *phone.PhoneNumber,
		PhonePoolId: phone.ID.String(),
	})
	if err != nil {
		panic(err)
	}
	return out.Bytes()
}
