package model

import (
	"context"
	"fmt"

	"github.com/aqaurius6666/apiservice/src/internal/db/business"
	"github.com/aqaurius6666/apiservice/src/internal/db/contact"
	"github.com/aqaurius6666/apiservice/src/internal/db/user"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/apiservice/src/pb"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

var (
	_ UserModel = (*ServerModel)(nil)
)

type UserModel interface {
	ListUsers(context.Context, *user.Search) ([]*user.User, error)
	TotalUsers(context.Context, *user.Search) (*int64, error)
	GetUser(context.Context, *user.Search) (*user.User, error)

	GetCredential(ctx context.Context, data string) (id, pubb, priv string, isDef bool, err error)
	Register(ctx context.Context, mail, phone, pub, enpriv string, role c.ROLE, refCode string) (otpId string, err error)
	VerifyOTP(ctx context.Context, otpId, otp string) (err error)
	ChangeMail(ctx context.Context, id interface{}, mail string) (string, error)
	// ChangeMailOTP(ctx context.Context, otpId, mail string) (string, error)
	ChangePassword(ctx context.Context, id, pub, enc string) error
	ForgotPasswordOTP(ctx context.Context, otpId, pub, enc string) error
	ForgotPassword(ctx context.Context, mail string) (otpId string, err error)
	CreateUser(ctx context.Context, id interface{}, mail, phone string) (*user.User, error)
	CheckCredential(ctx context.Context, identifier string) (bool, error)
	ResendOtp(ctx context.Context, otpId string) error

	ChangePassAndMail(context context.Context, userId string, mail string, pub string, enc string) error

	GetUserById(ctx context.Context, id interface{}) (*user.User, error)
	UpdateUser(ctx context.Context, id interface{}, u *user.User) error

	GetRegistrationProccess(context.Context, interface{}) (c.REGISTRATION_PROCESS, error)
	GetRegistrationProccessForUser(context.Context, interface{}) (c.REGISTRATION_PROCESS, error)

	ConvertUserToProto(*user.User) *pb.User
	ConvertUsersToProtos([]*user.User) []*pb.User
}

func (s *ServerModel) ConvertUsersToProtos(u []*user.User) []*pb.User {
	arr := make([]*pb.User, 0)
	for _, uu := range u {
		arr = append(arr, s.ConvertUserToProto(uu))
	}
	return arr
}

func (s *ServerModel) ConvertUserToProto(u *user.User) *pb.User {
	upb := new(pb.User)
	if u.ID != uuid.Nil {
		upb.Id = u.ID.String()
	}
	if u.AvatarUrl != nil {
		upb.Image = *u.AvatarUrl
	}
	if u.Mail != nil {
		upb.Mail = *u.Mail
	}
	if u.Phone != nil {
		upb.Phone = *u.Phone
	}
	if u.ContactId != uuid.Nil {
		upb.ContactId = u.ContactId.String()
	}
	if u.LastName != nil {
		upb.LastName = *u.LastName
	}
	if u.FirstName != nil {
		upb.FirstName = *u.FirstName
	}
	if u.Status != nil {
		upb.Status = c.ACCOUNT_STATUS(*u.Status)
	}
	if u.FirstName != nil && u.LastName != nil {
		upb.Name = fmt.Sprintf("%s %s", *u.FirstName, *u.LastName)
	}
	if u.Contact != nil && u.Contact.Zipcode != nil {
		upb.Zipcode = *u.Contact.Zipcode
	}
	return upb
}

func (s *ServerModel) GetRegistrationProccess(ctx context.Context, id interface{}) (c.REGISTRATION_PROCESS, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetRegistrationProccess))
	defer span.End()

	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return 0, err
	}
	bus, err := s.Repo.SelectBusiness(ctx, &business.Search{
		DefaultSearchModel: database.DefaultSearchModel{
			Fields: []string{
				`"businesses"."name"`,
				`"businesses"."banner_url"`,
				`"businesses"."services"`,
				`"businesses"."logo_url"`,
				`"businesses"."phone"`,
				`"Contact"."state_id"`,
				`"Contact"."city"`,
				`"Contact"."address1"`,
				`"Contact"."address2"`,
				`"businesses"."zipcodes"`,
			},
		},
		Business: business.Business{BaseModel: database.BaseModel{ID: uid}},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return c.REGISTRATION_PROCESS_DONE, err
	}

	if bus.Name == nil || bus.LogoUrl == nil || bus.BannerUrl == nil || bus.Services == nil {
		return c.REGISTRATION_PROCESS_STEP_1, nil
	}
	if bus.Phone == nil || bus.Contact.StateId == uuid.Nil || bus.Contact.City == nil || bus.Contact.Address1 == nil || bus.Contact.Address2 == nil || bus.Contact.Zipcode == nil {
		return c.REGISTRATION_PROCESS_STEP_2, nil
	}
	if bus.Zipcodes == nil || len(bus.Zipcodes) == 0 {
		return c.REGISTRATION_PROCESS_STEP_3, nil
	}
	return c.REGISTRATION_PROCESS_DONE, nil
}

func (s *ServerModel) GetRegistrationProccessForUser(ctx context.Context, id interface{}) (c.REGISTRATION_PROCESS, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetRegistrationProccessForUser))
	defer span.End()

	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return 0, err
	}
	cus, err := s.Repo.SelectUserProcess(ctx, &user.Search{
		DefaultSearchModel: database.DefaultSearchModel{
			Fields: []string{
				`"users"."avatar_url"`,
				`"users"."first_name"`,
				`"users"."last_name"`,
				`"Contact"."state_id"`,
				`"Contact"."city"`,
				`"Contact"."address1"`,
				`"Contact"."zipcode"`,
			},
		},
		User: user.User{BaseModel: database.BaseModel{ID: uid}},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return c.REGISTRATION_PROCESS_DONE, err
	}

	if cus.FirstName == nil || cus.LastName == nil || cus.Contact.Address1 == nil || cus.Contact.StateId == uuid.Nil || cus.Contact.City == nil || cus.Contact.Zipcode == nil {
		return c.REGISTRATION_PROCESS_STEP_1, nil
	}

	return c.REGISTRATION_PROCESS_DONE, nil
}

func (s *ServerModel) CheckCredential(ctx context.Context, identifier string) (bool, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CheckCredential))
	defer span.End()

	existed, err := s.Auth.CheckCredential(ctx, identifier)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return false, err
	}

	return existed, nil
}

func (s *ServerModel) Register(ctx context.Context, mail, phone, pub, enpriv string, role c.ROLE, refCode string) (otpId string, err error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.Register))
	defer span.End()

	otpId, err = s.Auth.Register(ctx, mail, phone, pub, enpriv, role, refCode)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return "", err
	}
	return otpId, nil
}

func (s *ServerModel) GetUserById(ctx context.Context, id interface{}) (*user.User, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetUserById))
	defer span.End()

	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	u, err := s.Repo.SelectUser(ctx, &user.Search{
		User: user.User{
			BaseModel: database.BaseModel{ID: uid},
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return u, nil
}

func (s *ServerModel) GetUser(ctx context.Context, search *user.Search) (*user.User, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetUser))
	defer span.End()

	u, err := s.Repo.SelectUser(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return u, nil
}

func (s *ServerModel) ChangeMail(ctx context.Context, id interface{}, mail string) (string, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ChangeMail))
	defer span.End()

	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return "", err
	}
	otpId, err := s.Auth.ChangeMail(ctx, uid.String(), mail)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return "", err
	}
	return otpId, nil
}

func (s *ServerModel) UpdateUser(ctx context.Context, id interface{}, u *user.User) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UpdateUser))
	defer span.End()

	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	if err := s.Repo.UpdateUser(ctx, &user.Search{
		User: user.User{
			BaseModel: database.BaseModel{ID: uid},
		},
	}, u); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}

func (s *ServerModel) CreateUser(ctx context.Context, id interface{}, mail, phone string) (*user.User, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CreateUser))
	defer span.End()

	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	u, err := s.Repo.InsertUser(ctx, &user.User{
		BaseModel: database.BaseModel{
			ID: uid,
		},
		Mail:      &mail,
		Phone:     &phone,
		ContactId: uid,
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	_, err = s.Repo.InsertContact(ctx, &contact.Contact{
		BaseModel: database.BaseModel{ID: uid},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return u, nil
}

func (s *ServerModel) VerifyOTP(ctx context.Context, otpId, otp string) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.VerifyOTP))
	defer span.End()

	typ, err := s.Auth.VerifyOTP(ctx, otpId, otp)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	switch typ {
	case c.OTP_TYPE_REGISTER:
		return s.handleRegisterAfter(ctx, otpId)
	case c.OTP_TYPE_FORGOT_PASSWORD:
		return nil
	case c.OTP_TYPE_CHANGE_MAIL:
		return s.handleChangeMailAfter(ctx, otpId)
	case c.OTP_TYPE_CHANGE_MAIL_AND_PASS:
		return s.handleChangeMailAndPassAfter(ctx, otpId)
	}
	return nil
}

func (s *ServerModel) handleRegisterAfter(ctx context.Context, otpId string) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.handleRegisterAfter))
	defer span.End()
	id, mail, phone, role, refCode, err := s.Auth.RegisterOTP(ctx, otpId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	switch role {
	case c.ROLE_CUSTOMER:
		if _, err = s.CreateUser(ctx, id, mail, phone); err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err, ctx)
			return err
		}
	case c.ROLE_HANDYMAN:
		if _, err = s.CreateBusiness(ctx, id, mail, phone, refCode); err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err, ctx)
			return err
		}
		// if _, err = s.NewCustomer(ctx, id, mail); err != nil {
		// 	err = xerrors.Errorf("%w", err)
		// 	lib.RecordError(span, err, ctx)
		// 	return err
		// }
	}
	return nil
}

func (s *ServerModel) handleChangeMailAfter(ctx context.Context, otpId string) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.handleChangeMailAfter))
	defer span.End()

	id, mail, role, err := s.Auth.ChangeMailOTP(ctx, otpId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	switch role {
	case c.ROLE_CUSTOMER:
		if err := s.UpdateUser(ctx, id, &user.User{
			Mail: &mail,
		}); err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err, ctx)
			return err
		}
	case c.ROLE_HANDYMAN:
		if err := s.UpdateBusiness(ctx, id, &business.Business{
			Mail: &mail,
		}); err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err, ctx)
			return err
		}
	}
	return nil
}

func (s *ServerModel) handleChangeMailAndPassAfter(ctx context.Context, otpId string) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.handleChangeMailAndPassAfter))
	defer span.End()

	id, mail, role, err := s.Auth.ChangeMailAndPassOTP(ctx, otpId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	switch role {
	case c.ROLE_CUSTOMER:
		if err := s.UpdateUser(ctx, id, &user.User{
			Mail: &mail,
		}); err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err, ctx)
			return err
		}
	case c.ROLE_HANDYMAN:
		if err := s.UpdateBusiness(ctx, id, &business.Business{
			Mail: &mail,
		}); err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err, ctx)
			return err
		}
	}
	return nil
}

func (s *ServerModel) ListUsers(ctx context.Context, search *user.Search) ([]*user.User, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListUsers))
	defer span.End()

	u, err := s.Repo.ListUsers(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return u, nil
}

func (s *ServerModel) TotalUsers(ctx context.Context, search *user.Search) (*int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.TotalUsers))
	defer span.End()

	c, err := s.Repo.TotalUsers(ctx, search)

	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return c, nil
}

func (s *ServerModel) ForgotPassword(ctx context.Context, mail string) (otpId string, err error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ForgotPassword))
	defer span.End()

	otpId, err = s.Auth.ForgotPassword(ctx, mail)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return "", err
	}
	return otpId, nil
}

func (s *ServerModel) ForgotPasswordOTP(ctx context.Context, otpId, pub, enc string) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ForgotPasswordOTP))
	defer span.End()

	err := s.Auth.ForgotPasswordOTP(ctx, otpId, pub, enc)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}

func (s *ServerModel) ChangePassword(ctx context.Context, id, pub, enc string) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ChangePassword))
	defer span.End()

	err := s.Auth.ChangePassword(ctx, id, pub, enc)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}

func (s *ServerModel) GetCredential(ctx context.Context, data string) (id, pubb, priv string, isDef bool, err error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetCredential))
	defer span.End()

	id, pub, priv, isDef, err := s.Auth.GetCredential(ctx, data)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return "", "", "", false, err
	}
	return id, pub, priv, isDef, nil
}

func (s *ServerModel) ResendOtp(ctx context.Context, otpId string) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ResendOtp))
	defer span.End()

	err := s.Auth.ResendOTP(ctx, otpId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}

func (s *ServerModel) ChangePassAndMail(ctx context.Context, userId string, mail string, pub string, enc string) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ChangePassAndMail))
	defer span.End()

	err := s.Auth.ChangeMailAndPass(ctx, userId, mail, pub, enc)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}

	return nil
}
