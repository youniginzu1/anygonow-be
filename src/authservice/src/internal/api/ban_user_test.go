package api

import (
	"context"
	"testing"

	"github.com/aqaurius6666/authservice/src/pb/authpb"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestBanUser(t *testing.T) {
	TEST_CASE := []map[string]string{
		{
			"id": "010a10c7-a807-4d7f-9bd7-d082a66e1e01",
		},
	}
	ctx := context.Background()
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	)
	assert.Nil(t, err)
	client := authpb.NewAuthServiceClient(conn)

	for _, c := range TEST_CASE {
		res, err := client.BanUser(ctx, &authpb.BanUserRequest{
			Id:     c["id"],
			Status: true,
		},
		)
		if err != nil {
			stt, ok := status.FromError(err)
			assert.True(t, ok)
			t.Error(stt.Message())
		}
		t.Log(res)
	}
}
