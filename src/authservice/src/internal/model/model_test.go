package model

import (
	"context"
	"testing"

	"github.com/aqaurius6666/authservice/src/internal/db"
	"github.com/aqaurius6666/authservice/src/internal/db/user"
	"github.com/aqaurius6666/authservice/src/internal/var/c"
	"github.com/aqaurius6666/authservice/src/services/mailservice"
	"github.com/aqaurius6666/go-utils/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestShouldSend(t *testing.T) {
	logger := logrus.New()
	ctx := context.Background()
	model, err := ModelMock(ctx, logger, ModelMockOptions{
		DBDsn:    db.DBDsn("postgresql://root:root@localhost:50010/defaultdb?sslmode=disable"),
		MailAddr: mailservice.MailServiceAddr(""),
	})
	assert.Nil(t, err, err)
	testCase := []map[string]interface{}{
		{
			"mail":     "aqaurius6666@gmail.com",
			"type":     c.OTP_TYPE_CHANGE_MAIL,
			"expected": true,
			"err":      nil,
		},
	}

	for _, tc := range testCase {
		mail := tc["mail"].(string)
		typ := tc["type"].(c.OTP_TYPE)
		exp := tc["expected"].(bool)
		ok, err := model.ShouldSend(ctx, &mail, typ)
		if e, okk := tc["err"].(error); okk && e != nil {
			assert.Equal(t, e.Error(), err.Error())
		} else {
			assert.Equal(t, exp, ok)
		}
	}
}

func TestCreateRegisterOTP(t *testing.T) {
	logger := logrus.New()
	ctx := context.Background()
	model, err := ModelMock(ctx, logger, ModelMockOptions{
		DBDsn:    db.DBDsn("postgresql://root:root@localhost:50010/defaultdb?sslmode=disable"),
		MailAddr: mailservice.MailServiceAddr("mailservice:50051"),
	})
	assert.Nil(t, err, err)
	testCase := []map[string]interface{}{
		{
			"user": &user.User{
				Username:            utils.StrPtr("aqaurius6666+5@gmail.com"),
				PublicKey:           utils.StrPtr("asd"),
				EncryptedPrivateKey: utils.StrPtr("asd"),
				RoleID:              uuid.MustParse(""),
				Mail:                utils.StrPtr("aqaurius6666+5@gmail.com"),
				Phone:               utils.StrPtr("aqaurius6666+5@gmail.com"),
			},
			"expected": true,
			"err":      nil,
		},
	}

	for _, tc := range testCase {
		user := tc["user"].(*user.User)
		exp := tc["expected"].(bool)
		ok, err := model.CreateRegisterOTP(ctx, user)
		if e, okk := tc["err"].(error); okk && e != nil {
			assert.Equal(t, e.Error(), err.Error())
		} else {
			assert.Equal(t, exp, ok)
		}
	}
}
