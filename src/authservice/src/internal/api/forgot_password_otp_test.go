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

func TestForgotPasswordOTP(t *testing.T) {
	TEST_CASE := []map[string]string{
		{
			"otpId":               "e61e1e01-5a0f-46cc-a45c-eb7e534d851e",
			"publicKey":           "dasdadasd",
			"encryptedPrivateKey": "dasdadasdsad",
		},
	}
	ctx := context.Background()
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.Nil(t, err)
	client := authpb.NewAuthServiceClient(conn)

	for _, c := range TEST_CASE {
		res, err := client.ForgotPasswordOTP(ctx, &authpb.ForgotPasswordOTPRequest{
			OtpId:               c["otpId"],
			EncryptedPrivateKey: c["encryptedPrivateKey"],
			PublicKey:           c["publicKey"],
		})
		if err != nil {
			stt, ok := status.FromError(err)
			assert.True(t, ok)
			t.Error(stt.Message())
		}
		t.Log(res)
	}
}
