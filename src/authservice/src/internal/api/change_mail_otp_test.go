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

func TestChangeMailOTP(t *testing.T) {
	TEST_CASE := []map[string]string{
		{
			"otpId": "4c431a63-0f78-4388-af6e-c933dc478dac",
			"otp":   "172711",
		},
	}
	ctx := context.Background()
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.Nil(t, err)
	client := authpb.NewAuthServiceClient(conn)

	for _, c := range TEST_CASE {
		res, err := client.ChangeMailOTP(ctx, &authpb.ChangeMailOTPRequest{
			OtpId: c["otpId"],
		})
		if err != nil {
			stt, ok := status.FromError(err)
			assert.True(t, ok)
			t.Error(stt.Message())
		}
		t.Log(res)
	}
}
