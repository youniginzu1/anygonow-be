package twilloclient

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConversation(t *testing.T) {

	ctx := context.Background()
	client := MockClient(t)
	friendlyName := "test"
	participants := []string{"+17066755365"}
	err := client.NewConversation(ctx, &friendlyName, participants...)
	assert.Nil(t, err)
}

func TestDeleteConversation(t *testing.T) {
	ctx := context.Background()
	client := MockClient(t)
	convId := "CH0436610167a44c5384133d391efef851"
	err := client.DeleteConversation(ctx, convId)
	assert.Nil(t, err)
}

func TestSendMessage(t *testing.T) {
	ctx := context.Background()
	client := MockClient(t)
	message := "Hello"
	err := client.SendMessage(ctx, "+14706844698", "+17066755365", message)
	assert.Nil(t, err)
}

func TestListAvailablePhone(t *testing.T) {
	ctx := context.Background()
	client := MockClient(t)
	phone, err := client.ListAvailablePhoneNumber(ctx)
	assert.Nil(t, err)
	assert.NotEmpty(t, phone)
	t.Log(*phone)
}

func TestListResourcePhone(t *testing.T) {
	ctx := context.Background()
	client := MockClient(t)
	phones, sids, err := client.ListResourcePhone(ctx)
	assert.Nil(t, err)
	assert.NotEmpty(t, phones)
	assert.NotEmpty(t, sids)
	t.Log(phones)
	t.Log(sids)
}

func TestRemovePhoneNumber(t *testing.T) {
	ctx := context.Background()
	client := MockClient(t)
	sid := "PNd3b30f93c40d75af71a30ff71abcd99d"
	err := client.ReleasePhoneNumber(ctx, sid)
	assert.Nil(t, err)
}
