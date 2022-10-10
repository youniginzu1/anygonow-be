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

func (s *ApiServer) ForgotPasswordOTP(ctx context.Context, req *authpb.ForgotPasswordOTPRequest) (*authpb.ForgotPasswordOTPResponse, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ForgotPasswordOTP))
	defer span.End()
	o, err := s.Model.GetValidOtpById(ctx, req.OtpId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		panic(err)
	}
	if !utils.BoolVal(o.Commited) || c.OTP_TYPE(utils.IntVal(o.Type)) != c.OTP_TYPE_FORGOT_PASSWORD {
		err = xerrors.Errorf("%w", e.ErrOTPInvalid)
		lib.RecordError(span, err)
		panic(err)
	}
	if err = s.Model.UpdateKey(ctx, o.UserId, req.PublicKey, req.EncryptedPrivateKey); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		panic(err)
	}
	if err := s.Model.DeleteOTPById(ctx, o.ID); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		// panic(err) // SHOULDN'T PANIC HERE
	}
	return &authpb.ForgotPasswordOTPResponse{}, nil
}
