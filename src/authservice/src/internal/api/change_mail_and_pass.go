package api

import (
	"context"

	"github.com/aqaurius6666/authservice/src/internal/db/otp"
	"github.com/aqaurius6666/authservice/src/internal/db/user"
	"github.com/aqaurius6666/authservice/src/internal/lib"
	"github.com/aqaurius6666/authservice/src/internal/var/c"
	"github.com/aqaurius6666/authservice/src/pb/authpb"
	"github.com/aqaurius6666/go-utils/database"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

func (s *ApiServer) ChangeMailAndPass(ctx context.Context, req *authpb.ChangeMailAndPassRequest) (*authpb.ChangeMailAndPassResponse, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ChangeMailAndPass))
	defer span.End()
	u, err := s.Model.GetUserById(ctx, req.UserId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		panic(err)
	}
	if *u.Mail == req.Mail {
		return &authpb.ChangeMailAndPassResponse{}, nil
	}
	if _, err = s.Model.GetUserByMail(ctx, req.Mail); err == nil {
		err = xerrors.Errorf("%w", user.ErrEmailExisted)
		lib.RecordError(span, err)
		panic(err)
	}
	o, err := s.Model.CreateChangeMailAndPassOTP(ctx, req.UserId, &user.User{
		PublicKey:           &req.PublicKey,
		EncryptedPrivateKey: &req.EncryptedPrivateKey,
		Mail:                &req.Mail,
		BaseModel: database.BaseModel{
			ID: u.ID,
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		panic(err)
	}
	loopFunc := func(ctx context.Context, funcname string, o *otp.Otp) {
		ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, funcname)
		defer span.End()
		var err error
		limit := 4
		for i := 0; i < limit; i++ {
			if err = s.Model.SendOTP(ctx, o); err == nil {
				return
			}
			s.Logger.Error(err)
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
	}
	go loopFunc(context.Background(), lib.GetFunctionName(loopFunc), o)
	return &authpb.ChangeMailAndPassResponse{}, nil
}
