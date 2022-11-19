package api

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/aqaurius6666/mailservice/src/pb/mailpb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestSendNotification(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.Nil(t, err)
	client := mailpb.NewMailServiceClient(conn)
	id := uuid.New()
	message := map[string]string{
		"id":  id.String(),
		"seq": fmt.Sprint(id.ClockSequence()),
	}
	mess, err := json.Marshal(message)
	assert.Nil(t, err)
	res, err := client.SendNotification(ctx, &mailpb.SendNotificationRequest{
		To:      "11526c9b-b56a-4d79-9f7b-4e1162360f0f",
		Body:    "This test message",
		Message: string(mess),
	})
	assert.Nil(t, err)
	stt, _ := status.FromError(err)
	t.Log(stt.Message())
	t.Log(res)
}
