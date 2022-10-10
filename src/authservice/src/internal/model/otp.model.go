package model

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/aqaurius6666/authservice/src/internal/db/otp"
	"github.com/aqaurius6666/authservice/src/internal/db/user"
	"github.com/aqaurius6666/authservice/src/internal/lib"
	"github.com/aqaurius6666/authservice/src/internal/lib/template"
	"github.com/aqaurius6666/authservice/src/internal/var/c"
	"github.com/aqaurius6666/authservice/src/internal/var/e"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/aqaurius6666/go-utils/utils"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

type OtpModel interface {
	GetValidOtpById(ctx context.Context, id interface{}) (*otp.Otp, error)
	CreateRegisterOTP(ctx context.Context, u *user.User) (*otp.Otp, error)
	CreateForgotOTP(ctx context.Context, u *user.User) (*otp.Otp, error)
	CreateChangeMailOTP(ctx context.Context, id interface{}, u *user.User) (*otp.Otp, error)
	CreateChangeMailAndPassOTP(ctx context.Context, id interface{}, u *user.User) (*otp.Otp, error)
	GenCode() *string
	SendOTP(ctx context.Context, o *otp.Otp) error
	ResendOTP(ctx context.Context, o *otp.Otp) error
	IsMatch(ctx context.Context, expected *otp.Otp, actual string) bool
	DeleteOTPById(ctx context.Context, id interface{}) error
	HandleRegisterMetadata(ctx context.Context, mt string) (*user.User, error)
	ChangeMailAndPassMetadata(ctx context.Context, mt string) (*user.User, error)
	CommitOTP(ctx context.Context, o *otp.Otp) error
	ShouldSend(ctx context.Context, mail *string, typ c.OTP_TYPE) (bool, error)

	genTemplateEmail(o *otp.Otp) ([]byte, error)
	createOTP(ctx context.Context, id uuid.UUID, _type c.OTP_TYPE, metadata *string, mail *string) (*otp.Otp, error)
}

func (s *ServerModel) ShouldSend(ctx context.Context, mail *string, typ c.OTP_TYPE) (bool, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ShouldSend))
	defer span.End()
	return true, nil
	o, err := s.Repo.SelectOtp(ctx, &otp.Search{
		Otp: otp.Otp{
			Type:     utils.IntPtr(int(typ)),
			Mail:     mail,
			Commited: utils.BoolPtr(false),
		},
	})
	if err != nil {
		return true, nil
	}
	if o.ExpireTime < time.Now().Unix() {
		return true, nil
	}

	return false, e.ErrOTPSpam
}

func (s *ServerModel) ResendOTP(ctx context.Context, o *otp.Otp) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ResendOTP))
	defer span.End()
	o.Code = s.GenCode()
	o.ExpireTime = time.Now().Add(c.OTP_EXPIRE_TIME).UnixMilli()
	bz, err := s.genTemplateEmail(o)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	if err := s.Mail.SendMail(ctx, *o.Mail, bz); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	if err := s.Repo.UpdateOTP(ctx, &otp.Search{
		Otp: otp.Otp{BaseModel: database.BaseModel{ID: o.ID}},
	}, o); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	return nil
}

func (s *ServerModel) CreateChangeMailOTP(ctx context.Context, id interface{}, u *user.User) (*otp.Otp, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CreateForgotOTP))
	defer span.End()
	bz, err := json.Marshal(u)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	mt := string(bz)
	return s.createOTP(ctx, u.ID, c.OTP_TYPE_CHANGE_MAIL, &mt, u.Mail)
}

func (s *ServerModel) CreateChangeMailAndPassOTP(ctx context.Context, id interface{}, u *user.User) (*otp.Otp, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CreateChangeMailAndPassOTP))
	defer span.End()
	bz, err := json.Marshal(u)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	mt := string(bz)
	return s.createOTP(ctx, u.ID, c.OTP_TYPE_CHANGE_MAIL_AND_PASS, &mt, u.Mail)
}

func (s *ServerModel) CreateForgotOTP(ctx context.Context, u *user.User) (*otp.Otp, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CreateForgotOTP))
	defer span.End()
	bz, err := json.Marshal(u)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	mt := string(bz)
	return s.createOTP(ctx, u.ID, c.OTP_TYPE_FORGOT_PASSWORD, &mt, u.Mail)
}

func (s *ServerModel) DeleteOTPById(ctx context.Context, id interface{}) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.DeleteOTPById))
	defer span.End()
	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}

	if err := s.Repo.DeleteOTP(ctx, &otp.Search{
		Otp: otp.Otp{
			BaseModel: database.BaseModel{
				ID: uid,
			},
		},
	}); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	return nil
}

func (s *ServerModel) HandleRegisterMetadata(ctx context.Context, mt string) (*user.User, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.HandleRegisterMetadata))
	defer span.End()
	var usr user.User
	if err := json.Unmarshal([]byte(mt), &usr); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	u, err := s.CreateUser(ctx, &user.User{
		Username:            usr.Username,
		PublicKey:           usr.PublicKey,
		EncryptedPrivateKey: usr.EncryptedPrivateKey,
		Mail:                usr.Mail,
		Phone:               usr.Phone,
		RoleID:              usr.RoleID,
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	u.RefCode = utils.StrPtr(utils.StrVal(usr.RefCode))
	return u, nil
}

func (s *ServerModel) ChangeMailAndPassMetadata(ctx context.Context, mt string) (*user.User, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ChangeMailAndPassMetadata))
	defer span.End()
	var usr user.User
	if err := json.Unmarshal([]byte(mt), &usr); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	return &usr, nil
}

func (s *ServerModel) IsMatch(ctx context.Context, expected *otp.Otp, actual string) bool {
	_, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.IsMatch))
	defer span.End()
	matched := *expected.Code == actual
	return matched
}

func (s *ServerModel) CommitOTP(ctx context.Context, o *otp.Otp) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CommitOTP))
	defer span.End()
	if err := s.Repo.UpdateOTP(ctx, &otp.Search{
		Otp: otp.Otp{
			BaseModel: database.BaseModel{ID: o.ID},
		},
	}, &otp.Otp{
		Commited: utils.BoolPtr(true),
	}); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	return nil
}

func (s *ServerModel) GetValidOtpById(ctx context.Context, id interface{}) (*otp.Otp, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetValidOtpById))
	defer span.End()
	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	o, err := s.Repo.SelectOtp(ctx, &otp.Search{
		Otp: otp.Otp{BaseModel: database.BaseModel{
			ID: uid,
		}},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	if o.ExpireTime < time.Now().UnixMilli() {
		// go s.DeleteOTPById(o.ID)
		err = xerrors.Errorf("%w", e.ErrOTPExpired)
		lib.RecordError(span, err)
		return nil, err
	}
	return o, nil
}
func (s *ServerModel) createOTP(ctx context.Context, id uuid.UUID, _type c.OTP_TYPE, metadata *string, mail *string) (*otp.Otp, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.createOTP))
	defer span.End()
	if ok, err := s.ShouldSend(ctx, mail, _type); err != nil || !ok {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	o, err := s.Repo.InsertOtp(ctx, &otp.Otp{
		Type:       utils.IntPtr(int(_type)),
		UserId:     id,
		ExpireTime: time.Now().Add(c.OTP_REGISTER_TIME).UnixMilli(),
		Code:       s.GenCode(),
		Metadata:   metadata,
		Mail:       mail,
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	return o, nil
}

func (s *ServerModel) SendOTP(ctx context.Context, o *otp.Otp) error {
	_, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.SendOTP))
	defer span.End()
	bz, err := s.genTemplateEmail(o)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
	}
	err = s.Mail.SendMail(ctx, *o.Mail, bz)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
	}
	return nil
}

func (s *ServerModel) CreateRegisterOTP(ctx context.Context, u *user.User) (*otp.Otp, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CreateRegisterOTP))
	defer span.End()
	bz, err := json.Marshal(u)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	mt := string(bz)
	return s.createOTP(ctx, uuid.Nil, c.OTP_TYPE_REGISTER, &mt, u.Mail)
}

func (s *ServerModel) GenCode() *string {
	rand.Seed(time.Now().Unix())
	r := rand.Int() % 999999
	str := fmt.Sprintf("%06d", r)
	return &str
}

func (s *ServerModel) genTemplateEmail(o *otp.Otp) ([]byte, error) {
	typ := c.OTP_TYPE(utils.IntVal(o.Type))
	switch typ {
	case c.OTP_TYPE_FORGOT_PASSWORD:
		return template.ForgotMailTemplate(*o.Mail, *o.Code, o.ID.String()), nil
	case c.OTP_TYPE_REGISTER:
		return template.RegisterMailTemplate(*o.Mail, *o.Code, o.ID.String()), nil
	case c.OTP_TYPE_CHANGE_MAIL_AND_PASS:
		return template.ChangeMailAndPassTemplate(*o.Mail, *o.Code, o.ID.String()), nil
	default:
		return nil, xerrors.Errorf("invalid type")
	}
}
