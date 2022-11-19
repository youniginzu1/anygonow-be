package chatservice

import (
	"context"
	"testing"

	"github.com/aqaurius6666/apiservice/src/pb/chatpb"
	"github.com/stretchr/testify/assert"
)

func TestChatservice(t *testing.T) {
	ctx := context.Background()
	chat, err := ConnectClient(ctx, ChatserviceAddr("localhost:50050"))
	assert.Nil(t, err, err)
	conv, err := chat.NewConversation(ctx, &chatpb.NewConversationRequest{
		MemberIds: []string{"5c5db1ba-a7a1-4dc5-b657-93ca90331126", "2472d3cb-68bc-4ef8-9d1c-6bf02f87dadc"},
	})
	assert.Nil(t, err, err)
	t.Log(conv)
}
