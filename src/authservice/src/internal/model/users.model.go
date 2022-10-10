package model

import (
	"context"
	"fmt"

	"github.com/aqaurius6666/authservice/src/internal/db/user"
	"github.com/aqaurius6666/authservice/src/internal/lib"
	"github.com/aqaurius6666/authservice/src/internal/var/c"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/aqaurius6666/go-utils/utils"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

type UserModel interface {
	ListUsers(context.Context, *user.Search) ([]*user.User, error)
	GetUserByIdPk(ctx context.Context, id interface{}, pubkey string) (*user.User, error)
	GetUserById(ctx context.Context, id interface{}) (*user.User, error)
	GetUserByUsername(ctx context.Context, username string) (*user.User, error)
	GetUserByPhone(ctx context.Context, phone string) (*user.User, error)
	GetUserByMail(ctx context.Context, email string) (*user.User, error)
	UpdateKey(ctx context.Context, id interface{}, pubkey, epk string) error
	UpdateMail(ctx context.Context, id interface{}, mail string) error
	CreateUser(ctx context.Context, u *user.User) (*user.User, error)
	CheckUserExisted(ctx context.Context, u *user.User) error
	BanUser(ctx context.Context, id interface{}, status bool) error
	DeleteUser(ctx context.Context, id interface{}) error
	UpdateMailAndKey(ctx context.Context, id interface{}, mail, pubkey, epk string) error
}

func (s *ServerModel) DeleteUser(ctx context.Context, id interface{}) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.DeleteUser))
	defer span.End()

	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	if err := s.Repo.DeleteUser(ctx, &user.Search{
		User: user.User{
			BaseModel: database.BaseModel{ID: uid},
		},
	}); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	return nil
}
func (s *ServerModel) BanUser(ctx context.Context, id interface{}, status bool) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.BanUser))
	defer span.End()

	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	if err := s.Repo.UpdateUser(ctx, &user.Search{
		User: user.User{
			BaseModel: database.BaseModel{ID: uid},
		},
	}, &user.User{
		IsActive: utils.BoolPtr(status),
	}); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	return nil
}

func (s *ServerModel) UpdateMail(ctx context.Context, id interface{}, mail string) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UpdateMail))
	defer span.End()

	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	if err := s.Repo.UpdateUser(ctx, &user.Search{
		User: user.User{
			BaseModel: database.BaseModel{ID: uid},
		},
	}, &user.User{
		Mail: &mail,
	}); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	return nil
}

func (s *ServerModel) UpdateMailAndKey(ctx context.Context, id interface{}, mail, pubkey, epk string) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UpdateMailAndKey))
	defer span.End()

	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	if err := s.Repo.UpdateUser(ctx, &user.Search{
		User: user.User{
			BaseModel: database.BaseModel{ID: uid},
		},
	}, &user.User{
		Mail:                &mail,
		PublicKey:           &pubkey,
		EncryptedPrivateKey: &epk,
	}); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	return nil
}

func (s *ServerModel) CheckUserExisted(ctx context.Context, u *user.User) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CheckUserExisted))
	defer span.End()

	var err error
	if _, err = s.Repo.SelectUser(ctx, &user.Search{
		User: user.User{
			Mail: u.Mail,
		},
	}); err == nil {
		err = xerrors.Errorf("%w", user.ErrEmailExisted)
		lib.RecordError(span, err)
		return err
	}

	if _, err = s.Repo.SelectUser(ctx, &user.Search{
		User: user.User{
			Phone: u.Phone,
		},
	}); err == nil {
		err = xerrors.Errorf("%w", user.ErrPhoneExisted)
		lib.RecordError(span, err)
		return err
	}
	return nil
}

func (s *ServerModel) CreateUser(ctx context.Context, u *user.User) (*user.User, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CreateUser))
	defer span.End()

	var err error
	u, err = s.Repo.InsertUser(ctx, u)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	return u, err
}

func (s *ServerModel) ListUsers(ctx context.Context, search *user.Search) ([]*user.User, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListUsers))
	defer span.End()

	u, err := s.Repo.ListUsers(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	return u, nil
}

func (s *ServerModel) GetUserById(ctx context.Context, id interface{}) (*user.User, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetUserById))
	defer span.End()
	fmt.Printf("id: %v\n", id)
	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	usr, err := s.Repo.SelectUser(ctx, &user.Search{
		User: user.User{
			BaseModel: database.BaseModel{
				ID: uid,
			},
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err

	}

	return usr, nil
}

func (s *ServerModel) GetUserByIdPk(ctx context.Context, id interface{}, pubkey string) (*user.User, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetUserByIdPk))
	defer span.End()

	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	usr, err := s.Repo.SelectUser(ctx, &user.Search{
		User: user.User{
			BaseModel: database.BaseModel{
				ID: uid,
			},
			PublicKey: &pubkey,
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}

	return usr, nil
}

func (s *ServerModel) GetUserByUsername(ctx context.Context, username string) (*user.User, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetUserByUsername))
	defer span.End()

	usr, err := s.Repo.SelectUser(ctx, &user.Search{
		User: user.User{
			Username: &username,
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	return usr, nil
}

func (s *ServerModel) GetUserByPhone(ctx context.Context, phone string) (*user.User, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetUserByPhone))
	defer span.End()

	usr, err := s.Repo.SelectUser(ctx, &user.Search{
		User: user.User{
			Phone: &phone,
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}

	return usr, nil
}

func (s *ServerModel) GetUserByMail(ctx context.Context, mail string) (*user.User, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetUserByMail))
	defer span.End()

	usr, err := s.Repo.SelectUser(ctx, &user.Search{
		User: user.User{
			Mail: &mail,
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}

	return usr, nil
}

func (s *ServerModel) UpdateKey(ctx context.Context, id interface{}, pubkey, epk string) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UpdateKey))
	defer span.End()

	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	err = s.Repo.UpdateUser(ctx, &user.Search{
		User: user.User{
			BaseModel: database.BaseModel{
				ID: uid,
			},
		},
	}, &user.User{
		PublicKey:           &pubkey,
		EncryptedPrivateKey: &epk,
		IsDefaultPassword:   utils.BoolPtr(false),
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}

	return nil
}
