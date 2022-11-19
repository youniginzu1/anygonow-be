package twilloclient

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/wire"
)

var MockSet = wire.NewSet(MockClient)

func MockClient(t *testing.T) Twilio {
	ctrl := gomock.NewController(t)
	return NewMockTwilio(ctrl)
}
