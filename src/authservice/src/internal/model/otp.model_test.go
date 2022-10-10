package model

import (
	"context"
	"testing"

	"github.com/aqaurius6666/authservice/src/internal/db/otp"
	"github.com/aqaurius6666/authservice/src/services/mailservice"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/aqaurius6666/go-utils/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestServerModel_SendOTP(t *testing.T) {
	ctx := context.Background()
	client, err := mailservice.ConnectClient(ctx, "localhost:50052")
	assert.Nil(t, err)
	mailClient := mailservice.ServiceGRPC{Ctx: ctx, Client: client}
	type fields struct {
		Ctx    context.Context
		Logger *logrus.Logger
		Mail   mailservice.Service
	}
	type args struct {
		ctx context.Context
		o   *otp.Otp
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "send otp",
			fields: fields{
				Ctx:    ctx,
				Logger: logrus.New(),
				Mail:   mailClient,
			},
			args: args{
				ctx: ctx,
				o: &otp.Otp{
					BaseModel: database.BaseModel{
						ID: uuid.New(),
					},
					UserId: uuid.New(),
					Type:   utils.IntPtr(1),
					Code:   utils.StrPtr("123456"),
					Mail:   utils.StrPtr("aqaurius1910+12@gmail.com"),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ServerModel{
				Ctx:    tt.fields.Ctx,
				Logger: tt.fields.Logger,
				Mail:   tt.fields.Mail,
			}
			if err := s.SendOTP(tt.args.ctx, tt.args.o); (err != nil) != tt.wantErr {
				t.Errorf("ServerModel.SendOTP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
