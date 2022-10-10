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

func (s *ApiServer) RegisterOTP(ctx context.Context, req *authpb.RegisterOTPRequest) (*authpb.RegisterOTPResponse, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.RegisterOTP))
	defer span.End()
	o, err := s.Model.GetValidOtpById(ctx, req.OtpId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		panic(err)
	}
	if !utils.BoolVal(o.Commited) || c.OTP_TYPE(utils.IntVal(o.Type)) != c.OTP_TYPE_REGISTER {
		err = xerrors.Errorf("%w", e.ErrOTPInvalid)
		lib.RecordError(span, err)
		panic(err)
	}
	u, err := s.Model.HandleRegisterMetadata(ctx, *o.Metadata)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		panic(err)
	}
	refCode := utils.StrPtr(utils.StrVal(u.RefCode))
	if err := s.Model.DeleteOTPById(ctx, o.ID); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
	}
	u, err = s.Model.GetUserById(ctx, u.ID)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		panic(err)
	}
	return &authpb.RegisterOTPResponse{
		Id:      u.ID.String(),
		Mail:    *u.Mail,
		Phone:   *u.Phone,
		Role:    c.ROLE(utils.IntVal(u.Role.Code)),
		RefCode: *refCode,
	}, nil
}
