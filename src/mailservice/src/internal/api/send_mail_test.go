package api

import (
	"context"
	"testing"

	"github.com/aqaurius6666/mailservice/src/pb/mailpb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestSendMail(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.Nil(t, err)
	client := mailpb.NewMailServiceClient(conn)

	res, err := client.SendMail(ctx, &mailpb.SendMailRequest{
		To:  "aqaurius1910+1@gmail.com",
		Msg: []byte("hello im from vn, please go to localhost/3000"),
	})
	assert.Nil(t, err)
	stt, _ := status.FromError(err)
	t.Log(stt.Message())
	t.Log(res)
}
