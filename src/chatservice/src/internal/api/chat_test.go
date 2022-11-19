package api

import (
	"context"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/aqaurius6666/chatservice/src/pb/chatpb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func newClient(ctx context.Context) (chatpb.ChatServiceClient, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, "localhost:50050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := chatpb.NewChatServiceClient(conn)
	return client, err
}

func TestChat(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := newClient(ctx)
	assert.Nil(t, err)
	stream, err := client.Chat(ctx)
	assert.Nil(t, err)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				return
			}
			t.Log(msg)
			if err != nil {
				return
			}
			if err := stream.Send(msg); err != nil {
				return
			}
		}
	}()
	wg.Wait()
}

func TestNewConversation(t *testing.T) {
	ctx := context.Background()
	client, err := newClient(ctx)
	assert.Nil(t, err)
	orderId := uuid.New().String()
	serviceId := "25483915-a70c-4161-bb15-617e0c3ddf3b"
	out, err := client.NewConversation(ctx, &chatpb.NewConversationRequest{
		ServiceId:    serviceId,
		OrderId:      orderId,
		MemberIds:    []string{"25483915-a70c-4161-bb15-617e0c3ddf3d"	, "61f2c92a-af19-4d78-b507-58fa3210d522"},
		PhoneNumbers: []string{"2604680522", "7066755368"},
	})
	assert.Nil(t, err)
	t.Log(err)
	t.Log(out)
	t.Log(orderId)
	t.Log(serviceId)
}

func TestCloseConversation(t *testing.T) {
	ctx := context.Background()
	client, err := newClient(ctx)
	assert.Nil(t, err)
	orderId := "0683ff9e-fdf2-4925-8c3b-308fc6971d2e"

	out, err := client.CloseConversation(ctx, &chatpb.CloseConversationRequest{
		OrderId: orderId,
	})
	assert.Nil(t, err)
	t.Log(err)
	t.Log(out)
	t.Log(orderId)
}
