package model

import (
	"context"

	"github.com/aqaurius6666/chatservice/src/pb/chatpb"
)

type Publisher interface {
	Publish(*chatpb.ChatMessage) error
}

type PublishV1 struct {
	ctx   context.Context
	topic string
}

func (s *PublishV1) Publish(msg *chatpb.ChatMessage) error {
	return nil
}

func NewPublisher(ctx context.Context, topic string) (Publisher, error) {
	return &PublishV1{
		ctx: ctx, topic: topic,
	}, nil
}
