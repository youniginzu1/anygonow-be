package model

import (
	"context"

	"github.com/aqaurius6666/chatservice/src/pb/chatpb"
)

type Subscriber interface {
	Subscribe(string) <-chan *chatpb.ChatMessage
}

type SubscribeV1 struct {
	ctx context.Context
}

func NewSubscriber(ctx context.Context) (Subscriber, error) {
	return &SubscribeV1{
		ctx: ctx,
	}, nil
}

func (s *SubscribeV1) Subscribe(topic string) <-chan *chatpb.ChatMessage {
	return nil
}
