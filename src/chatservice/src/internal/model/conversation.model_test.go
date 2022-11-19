package model

import (
	"context"
	"testing"

	"github.com/aqaurius6666/chatservice/src/internal/db"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func GetMockModel(t *testing.T) Server {
	mock, _ := NewMockModel(context.Background(), logrus.New(), db.DBDsn("postgresql://root:root@localhost:50000/defaultdb?sslmode=disable"), t)
	return mock
}
func TestSyncPhonePool(t *testing.T) {
	ctx := context.Background()
	mockModel := GetMockModel(t)
	err := mockModel.SyncPhonePool(ctx)
	assert.Nil(t, err)
}
