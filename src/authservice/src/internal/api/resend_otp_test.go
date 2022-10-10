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

func TestResendOTP(t *testing.T) {
	TEST_CASE := []map[string]string{
		{
			"id": "08ad2c0d-850b-4781-9e17-68af168fe87c",
		},
	}
	ctx := context.Background()
	conn, err := grpc.Dial("localhost:50011", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.Nil(t, err)
	client := authpb.NewAuthServiceClient(conn)

	for _, c := range TEST_CASE {
		res, err := client.ResendOTP(ctx, &authpb.ResendOTPRequest{
			Id: c["id"],
		})
		if err != nil {
			stt, ok := status.FromError(err)
			assert.True(t, ok)
			t.Error(stt.Message())
		}
		t.Log(res)

	}
}
