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

func TestVerifyOTP(t *testing.T) {
	TEST_CASE := []map[string]string{
		{
			"otpId": "2ebfe364-bf4b-4120-bed5-c723e1c5a0cd",
			"otp":   "271830",
		},
	}
	ctx := context.Background()
	conn, err := grpc.Dial("localhost:50011", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.Nil(t, err)
	client := authpb.NewAuthServiceClient(conn)

	for _, c := range TEST_CASE {
		res, err := client.VerifyOTP(ctx, &authpb.VerifyOTPRequest{
			OtpId: c["otpId"],
			Otp:   c["otp"],
		})
		if err != nil {
			stt, ok := status.FromError(err)
			assert.True(t, ok)
			t.Error(stt.Message())
		}
		t.Log(res)

	}
}
