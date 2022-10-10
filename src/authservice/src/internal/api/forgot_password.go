package api

import (
	"context"

	"github.com/aqaurius6666/authservice/src/internal/db/otp"
	"github.com/aqaurius6666/authservice/src/internal/lib"
	"github.com/aqaurius6666/authservice/src/internal/var/c"
	"github.com/aqaurius6666/authservice/src/pb/authpb"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

func (s *ApiServer) ForgotPassword(ctx context.Context, req *authpb.ForgotPasswordRequest) (*authpb.ForgotPasswordResponse, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ForgotPassword))
	defer span.End()
	usr, err := s.Model.GetUserByMail(ctx, req.Mail)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		panic(err)
	}
	o, err := s.Model.CreateForgotOTP(ctx, usr)
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
	return &authpb.ForgotPasswordResponse{
		Id:    usr.ID.String(),
		OtpId: o.ID.String(),
	}, nil
}
