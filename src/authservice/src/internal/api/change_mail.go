package api

import (
	"context"

	"github.com/aqaurius6666/authservice/src/internal/db/user"
	"github.com/aqaurius6666/authservice/src/internal/lib"
	"github.com/aqaurius6666/authservice/src/internal/var/c"
	"github.com/aqaurius6666/authservice/src/pb/authpb"
	"github.com/aqaurius6666/go-utils/database"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

func (s *ApiServer) ChangeMail(ctx context.Context, req *authpb.ChangeMailRequest) (*authpb.ChangeMailResponse, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ChangeMail))
	defer span.End()
	u, err := s.Model.GetUserById(ctx, req.Id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		panic(err)
	}
	if *u.Mail == req.Mail {
		return &authpb.ChangeMailResponse{}, nil
	}
	if _, err = s.Model.GetUserByMail(ctx, req.Mail); err == nil {
		err = xerrors.Errorf("%w", user.ErrEmailExisted)
		lib.RecordError(span, err)
		panic(err)
	}
	o, err := s.Model.CreateChangeMailOTP(ctx, req.Id, &user.User{
		Mail: &req.Mail,
		BaseModel: database.BaseModel{
			ID: u.ID,
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		panic(err)
	}
	return &authpb.ChangeMailResponse{
		OtpId: o.ID.String(),
	}, nil
}
