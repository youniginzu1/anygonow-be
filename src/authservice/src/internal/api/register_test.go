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

func TestRegisterNoOTP(t *testing.T) {
	TEST_CASE := []map[string]string{
		{
			"username":            "test",
			"publicKey":           "pk1",
			"encryptedPrivateKey": "pk",
			"mail":                "a@gma.com",
		},
	}
	ctx := context.Background()
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.Nil(t, err)
	client := authpb.NewAuthServiceClient(conn)

	for _, c := range TEST_CASE {
		res, err := client.RegisterNoOTP(ctx, &authpb.RegisterNoOTPRequest{
			Username:            c["username"],
			PublicKey:           c["publicKey"],
			EncryptedPrivateKey: c["encryptedPrivateKey"],
			Mail:                c["mail"],
		})
		if err != nil {
			stt, ok := status.FromError(err)
			assert.True(t, ok)
			t.Error(stt.Message())
		}
		t.Log(res)

	}
}

func TestRegisterOTP(t *testing.T) {
	TEST_CASE := []map[string]string{
		{
			"username":            "test+34",
			"publicKey":           "pk1+4",
			"encryptedPrivateKey": "pk+3",
			"mail":                "aqaurius1910@gmail.com",
			"phone":               "09011123",
		},
	}
	ctx := context.Background()
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.Nil(t, err)
	client := authpb.NewAuthServiceClient(conn)

	for _, c := range TEST_CASE {
		res, err := client.Register(ctx, &authpb.RegisterRequest{
			Username:            c["username"],
			PublicKey:           c["publicKey"],
			EncryptedPrivateKey: c["encryptedPrivateKey"],
			Mail:                c["mail"],
			Phone:               c["phone"],
		})
		if err != nil {
			stt, ok := status.FromError(err)
			assert.True(t, ok)
			t.Error(stt.Message())
		}
		t.Log(res)

	}
}
