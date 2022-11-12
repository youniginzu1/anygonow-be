package sms

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"golang.org/x/xerrors"
)

type SNSService struct {
	client *sns.SNS
}

func NewSNSService() (*SNSService, error) {
	sess, err := session.NewSession(aws.NewConfig().WithCredentials(credentials.NewChainCredentials([]credentials.Provider{
		&credentials.EnvProvider{},
	})))
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}

	client := sns.New(sess)
	return &SNSService{
		client: client,
	}, nil
}

func (s *SNSService) SendSMS(ctx context.Context, to string, msg string) error {
	out, err := s.client.PublishWithContext(ctx, &sns.PublishInput{
		Message:     aws.String(msg),
		PhoneNumber: aws.String(to),
	})
	if err != nil {
		return xerrors.Errorf("%w", err)
	}
	fmt.Printf("out.MessageId: %v\n", *out.MessageId)
	return nil
}
