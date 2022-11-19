package twilloclient

import (
	"context"
	"fmt"
	"os"

	"github.com/aqaurius6666/chatservice/src/internal/lib"
	"github.com/aqaurius6666/chatservice/src/internal/var/c"
	"github.com/aqaurius6666/go-utils/utils"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"github.com/twilio/twilio-go"
	"github.com/twilio/twilio-go/client"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	conversationDto "github.com/twilio/twilio-go/rest/conversations/v1"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

var Set = wire.NewSet(wire.Struct(new(TwilloClient), "*"), NewRestClient, NewValidator, wire.Bind(new(Twilio), new(*TwilloClient)))
var PROXY_ADDRESS = "+14706844698"

type TWILLO_SMS_CALLBACK_URL string
type Twilio interface {
	NewConversation(ctx context.Context, friendlyName *string, participant ...string) error
	DeleteConversation(ctx context.Context, convId string) error
	SendMessage(ctx context.Context, from, to string, message string) error
	ListAvailablePhoneNumber(ctx context.Context) (*string, error)
	ReleasePhoneNumber(ctx context.Context, sid string) error
	ListResourcePhone(ctx context.Context) (phones []string, sids []string, err error)
	BuyPhoneNumber(ctx context.Context, phone *string) (*string, *string, error)
}
type TwilloClient struct {
	Logger      *logrus.Logger
	Client      *twilio.RestClient
	CallbackUrl TWILLO_SMS_CALLBACK_URL
	Validator   client.RequestValidator
}

func NewRestClient() *twilio.RestClient {
	client := twilio.NewRestClient()
	return client
}

func NewValidator() client.RequestValidator {
	secret := os.Getenv("TWILIO_AUTH_TOKEN")
	return client.NewRequestValidator(secret)
}

// Unused
func (s *TwilloClient) NewConversation(ctx context.Context, friendlyName *string, participant ...string) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.NewConversation))
	defer span.End()

	conv, err := s.Client.ConversationsV1.CreateConversation(&conversationDto.CreateConversationParams{
		FriendlyName: friendlyName,
	})
	if err != nil {
		return err
	}
	for _, participant := range participant {
		participantRes, err := s.Client.ConversationsV1.CreateConversationParticipant(*conv.Sid, &conversationDto.CreateConversationParticipantParams{
			MessagingBindingAddress:      &participant,
			MessagingBindingProxyAddress: &PROXY_ADDRESS,
		})
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err)
			return err
		}
		fmt.Printf("participantRes.Sid: %v\n", *participantRes.Sid)
	}

	fmt.Printf("conv.Sid: %v\n", *conv.Sid)
	return nil
}

func (s *TwilloClient) DeleteConversation(ctx context.Context, convId string) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.NewConversation))
	defer span.End()

	err := s.Client.ConversationsV1.DeleteConversation(convId, &conversationDto.DeleteConversationParams{})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	return nil
}

func (s *TwilloClient) SendMessage(ctx context.Context, from, to string, message string) error {
	_, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.SendMessage))
	defer span.End()

	_, err := s.Client.Api.CreateMessage(&openapi.CreateMessageParams{
		Body: &message,
		From: &from,
		To:   &to,
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}

	return nil
}

func (s *TwilloClient) ListAvailablePhoneNumber(ctx context.Context) (*string, error) {
	_, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListAvailablePhoneNumber))
	defer span.End()

	out, err := s.Client.Api.ListAvailablePhoneNumberLocal("US", &openapi.ListAvailablePhoneNumberLocalParams{
		SmsEnabled: utils.BoolPtr(true),
		Limit:      utils.IntPtr(1),
	})
	if err != nil {
		return nil, err
	}

	if len(out) == 0 {
		return nil, xerrors.New("no available phone number")
	}
	return out[0].PhoneNumber, nil
}

func (s *TwilloClient) BuyPhoneNumber(ctx context.Context, phone *string) (phoneNumber *string, sid *string, err error) {
	_, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.BuyPhoneNumber))
	defer span.End()

	out, err := s.Client.Api.CreateIncomingPhoneNumberLocal(&openapi.CreateIncomingPhoneNumberLocalParams{
		PhoneNumber: phone,
		SmsMethod:   utils.StrPtr("POST"),
		SmsUrl:      utils.StrPtr(string(s.CallbackUrl)),
	})
	if err != nil {
		return nil, nil, err
	}

	return out.PhoneNumber, out.Sid, nil
}

func (s *TwilloClient) ListResourcePhone(ctx context.Context) (phones []string, sids []string, err error) {

	_, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListResourcePhone))
	defer span.End()

	out, err := s.Client.Api.ListIncomingPhoneNumber(&openapi.ListIncomingPhoneNumberParams{})
	if err != nil {
		return nil, nil, err
	}
	phones = make([]string, 0)
	sids = make([]string, 0)
	for _, phone := range out {
		phones = append(phones, *phone.PhoneNumber)
		sids = append(sids, *phone.Sid)
	}
	return phones, sids, nil
}

func (s *TwilloClient) ReleasePhoneNumber(ctx context.Context, sid string) error {
	_, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ReleasePhoneNumber))
	defer span.End()

	err := s.Client.Api.DeleteIncomingPhoneNumber(sid, &openapi.DeleteIncomingPhoneNumberParams{})
	if err != nil {
		return err
	}
	return nil
}

// Verify sinature
func (s *TwilloClient) VerifySignature(signature string, body []byte) bool {
	return s.Validator.ValidateBody(string(s.CallbackUrl), body, signature)
}
