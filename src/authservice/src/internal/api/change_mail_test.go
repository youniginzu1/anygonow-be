package api

import (
	"context"
	"testing"

	"github.com/aqaurius6666/authservice/src/pb/authpb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestChangeMail(t *testing.T) {
	TEST_CASE := []map[string]string{
		{
			"mail": "aqaurius19100@gmail.com",
			"id":   "868356c1-184e-4276-8784-4dd4efb12d5e",
		},
	}
	ctx := context.Background()
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.Nil(t, err)
	client := authpb.NewAuthServiceClient(conn)

	for _, c := range TEST_CASE {
		res, err := client.ChangeMail(ctx, &authpb.ChangeMailRequest{
			Id:   c["id"],
			Mail: c["mail"],
		})
		if err != nil {
			stt, ok := status.FromError(err)
			assert.True(t, ok)
			t.Error(stt.Message())
		}
		t.Log(res)
	}
}
