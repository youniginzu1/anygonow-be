package api

import (
	"context"
	"time"

	"github.com/aqaurius6666/authservice/src/internal/lib"
	"github.com/aqaurius6666/authservice/src/internal/var/c"
	"github.com/aqaurius6666/authservice/src/internal/var/e"
	"github.com/aqaurius6666/authservice/src/pb/authpb"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

func (s *ApiServer) ResendOTP(ctx context.Context, req *authpb.ResendOTPRequest) (*authpb.ResendOTPResponse, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ResendOTP))
	defer span.End()
	o, err := s.Model.GetValidOtpById(ctx, req.Id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		panic(err)
	}
	if o.UpdatedAt+c.OTP_SPAM_TIME.Milliseconds() > time.Now().UnixMilli() {
		err = xerrors.Errorf("%w", e.ErrOTPSpam)
		lib.RecordError(span, err)
		panic(err)
	}
	if err = s.Model.ResendOTP(ctx, o); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		panic(err)
	}
	return &authpb.ResendOTPResponse{}, nil
}
