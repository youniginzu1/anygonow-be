package model

import (
	"context"

	"github.com/aqaurius6666/apiservice/src/internal/db/business"
	"github.com/aqaurius6666/apiservice/src/internal/db/contact"
	"github.com/aqaurius6666/apiservice/src/internal/db/user"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/aqaurius6666/go-utils/utils"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

var (
	_ AdminModel = (*ServerModel)(nil)
)

type AdminModel interface {
	InactiveUser(ctx context.Context, id interface{}) (string, c.ROLE, error)
	InactiveBusiness(ctx context.Context, id interface{}) (string, c.ROLE, error)
	ActiveUser(ctx context.Context, id interface{}) (string, c.ROLE, error)
	ActiveBusiness(ctx context.Context, id interface{}) (string, c.ROLE, error)
	DeleteBusiness(ctx context.Context, id interface{}) error
	DeleteUser(ctx context.Context, id interface{}) error
}

func (s *ServerModel) DeleteUser(ctx context.Context, id interface{}) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.DeleteBusiness))
	defer span.End()

	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	_, _, err = s.Auth.DeleteUser(ctx, uid.String())
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	if err := s.Repo.DeleteUser(ctx, &user.Search{
		User: user.User{
			BaseModel: database.BaseModel{ID: uid},
		},
	}); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	if err := s.Repo.DeleteContact(ctx, &contact.Search{
		Contact: contact.Contact{
			BaseModel: database.BaseModel{ID: uid},
		},
	}); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}

func (s *ServerModel) DeleteBusiness(ctx context.Context, id interface{}) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.DeleteBusiness))
	defer span.End()

	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	_, _, err = s.Auth.DeleteUser(ctx, uid.String())
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	if err := s.Repo.DeleteBusiness(ctx, &business.Search{
		Business: business.Business{
			BaseModel: database.BaseModel{ID: uid},
		},
	}); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	if err := s.Repo.DeleteContact(ctx, &contact.Search{
		Contact: contact.Contact{
			BaseModel: database.BaseModel{ID: uid},
		},
	}); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}
func (s *ServerModel) InactiveUser(ctx context.Context, id interface{}) (string, c.ROLE, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.InactiveUser))
	defer span.End()

	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return "", 0, err
	}
	nid, r, err := s.Auth.BanUser(ctx, uid.String(), false)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return "", r, err
	}
	if err = s.UpdateUser(ctx, uid, &user.User{
		Status: utils.Int32Ptr(int32(c.ACCOUNT_STATUS_INACTIVE)),
	}); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return "", 0, err
	}
	return nid, r, nil
}

func (s *ServerModel) InactiveBusiness(ctx context.Context, id interface{}) (string, c.ROLE, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.InactiveBusiness))
	defer span.End()

	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return "", 0, err
	}
	nid, r, err := s.Auth.BanUser(ctx, uid.String(), false)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return "", 0, err
	}
	if err = s.UpdateBusiness(ctx, uid, &business.Business{
		Status: utils.Int32Ptr(int32(c.ACCOUNT_STATUS_INACTIVE)),
	}); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return "", 0, err
	}
	return nid, r, nil
}
func (s *ServerModel) ActiveUser(ctx context.Context, id interface{}) (string, c.ROLE, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ActiveUser))
	defer span.End()
	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return "", 0, err
	}
	nid, r, err := s.Auth.BanUser(ctx, uid.String(), true)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return "", 0, err
	}
	if err = s.UpdateUser(ctx, uid, &user.User{
		Status: utils.Int32Ptr(int32(c.ACCOUNT_STATUS_ACTIVE)),
	}); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return "", 0, err
	}
	return nid, r, nil
}

func (s *ServerModel) ActiveBusiness(ctx context.Context, id interface{}) (string, c.ROLE, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ActiveBusiness))
	defer span.End()

	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return "", 0, err
	}
	nid, r, err := s.Auth.BanUser(ctx, uid.String(), true)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return "", 0, err
	}
	if err = s.UpdateBusiness(ctx, uid, &business.Business{
		Status: utils.Int32Ptr(int32(c.ACCOUNT_STATUS_ACTIVE)),
	}); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return "", 0, err
	}
	return nid, r, nil
}
