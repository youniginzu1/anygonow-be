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

func TestCheckAuth(t *testing.T) {
	TEST_CASE := []map[string]string{
		{
			"header": `{"signature":"tDtVBhTTSnRSEdbKJMNXxgPhmxpd70Y9Nbrs85M40fAuNwPSXVCzubXYZilROMQPRTRYKjzLytQLbKUSlr0OKQ==","certificateInfo":{"id":"aef65976-d385-4525-add1-1cd2bcc6fd59","timestamp":1641721592445,"exp":2799360000000},"publicKey":"A9VMrb8olmifFj4QVhG63fIJDK1+kkKsdKE3bmm+E9Xx"}`,
			"method": `GET`,
			"body":   `{"data":{"id":"8475cd9a-7b6c-4354-aa39-cdb5c6ca762a","_actionType":"POST_V1-ADMINS-USERS-8475CD9A-7B6C-4354-AA39-CDB5C6CA762A-BAN","_timestamp":1641723014073},"_signature":"utBy5h2X+IqzEXfjo3rVp8bJxZN39MMzGVBCakcOhU4+cEx+vLskQU4TYCIYnLWgsi4X7Wiy0OBKL3tBD8TwsA=="}`,
		},
		{
			"header": `{"signature":"tDtVBhTTSnRSEdbKJMNXxgPhmxpd70Y9Nbrs85M40fAuNwPSXVCzubXYZilROMQPRTRYKjzLytQLbKUSlr0OKQ==","certificateInfo":{"id":"5ef65976-d385-4525-add1-1cd2bcc6fd59","timestamp":1641721592445,"exp":2799360000000},"publicKey":"A9VMrb8olmifFj4QVhG63fIJDK1+kkKsdKE3bmm+E9Xx"}`,
			"method": `POST`,
			"body":   `{"data":{"id":"8475cd9a-7b6c-4354-aa39-cdb5c6ca762a","_actionType":"POST_V1-ADMINS-USERS-8475CD9A-7B6C-4354-AA39-CDB5C6CA762A-BAN","_timestamp":1641723014073},"_signature":"utBy5h2X+IqzEXfjo3rVp8bJxZN39MMzGVBCakcOhU4+cEx+vLskQU4TYCIYnLWgsi4X7Wiy0OBKL3tBD8TwsA=="}`,
		},
		{},
	}
	ctx := context.Background()
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.Nil(t, err)
	client := authpb.NewAuthServiceClient(conn)

	for _, c := range TEST_CASE {
		res, err := client.CheckAuth(ctx, &authpb.CheckAuthRequest{
			Header: []byte(c["header"]),
			Method: c["method"],
			Body:   []byte(c["body"]),
		})
		stt, ok := status.FromError(err)
		assert.True(t, ok)
		t.Error(stt.Message())
		t.Log(res)
	}
}
