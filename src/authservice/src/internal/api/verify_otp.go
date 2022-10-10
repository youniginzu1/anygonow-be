package api

import (
	"context"

	"github.com/aqaurius6666/authservice/src/internal/lib"
	"github.com/aqaurius6666/authservice/src/internal/var/c"
	"github.com/aqaurius6666/authservice/src/internal/var/e"
	"github.com/aqaurius6666/authservice/src/pb/authpb"
	"github.com/aqaurius6666/go-utils/utils"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

func (s *ApiServer) VerifyOTP(ctx context.Context, req *authpb.VerifyOTPRequest) (*authpb.VerifyOTPResponse, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.VerifyOTP))
	defer span.End()
	o, err := s.Model.GetValidOtpById(ctx, req.OtpId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		panic(err)
	}
	if *o.Commited {
		return &authpb.VerifyOTPResponse{
			Ok:    true,
			Type:  c.OTP_TYPE(utils.IntVal(o.Type)),
			OtpId: o.ID.String(),
		}, nil
	}
	if match := s.Model.IsMatch(ctx, o, req.Otp); !match {
		err = xerrors.Errorf("%w", e.ErrOTPNotMatch)
		lib.RecordError(span, err)
		panic(err)
	}

	if err := s.Model.CommitOTP(ctx, o); err != nil {
		err = xerrors.Errorf("%w", e.ErrOTPNotMatch)
		lib.RecordError(span, err)
		panic(err)
	}
	return &authpb.VerifyOTPResponse{
		Ok:    true,
		Type:  c.OTP_TYPE(utils.IntVal(o.Type)),
		OtpId: o.ID.String(),
	}, nil
}
