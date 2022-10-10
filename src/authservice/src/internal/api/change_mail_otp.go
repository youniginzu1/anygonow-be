package api

import (
	"context"

	"github.com/aqaurius6666/authservice/src/internal/db/user"
	"github.com/aqaurius6666/authservice/src/internal/lib"
	"github.com/aqaurius6666/authservice/src/internal/var/c"
	"github.com/aqaurius6666/authservice/src/internal/var/e"
	"github.com/aqaurius6666/authservice/src/pb/authpb"
	"github.com/aqaurius6666/go-utils/utils"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

func (s *ApiServer) ChangeMailOTP(ctx context.Context, req *authpb.ChangeMailOTPRequest) (*authpb.ChangeMailOTPResponse, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ChangeMailOTP))
	defer span.End()
	o, err := s.Model.GetValidOtpById(ctx, req.OtpId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		panic(err)
	}
	if !utils.BoolVal(o.Commited) || c.OTP_TYPE(utils.IntVal(o.Type)) != c.OTP_TYPE_CHANGE_MAIL {
		err = xerrors.Errorf("%w", e.ErrOTPInvalid)
		lib.RecordError(span, err)
		panic(err)
	}
	if _, err = s.Model.GetUserByMail(ctx, *o.Mail); err == nil {
		err = xerrors.Errorf("%w", user.ErrEmailExisted)
		lib.RecordError(span, err)
		panic(err)
	}
	if err = s.Model.UpdateMail(ctx, o.UserId, *o.Mail); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		panic(err)
	}
	if err := s.Model.DeleteOTPById(ctx, o.ID); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		// panic(err) // SHOULDN'T PANIC HERE
	}
	u, err := s.Model.GetUserById(ctx, o.UserId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		panic(err)
	}

	return &authpb.ChangeMailOTPResponse{
		Id:   u.ID.String(),
		Mail: *u.Mail,
		Role: c.ROLE(utils.IntVal(u.Role.Code)),
	}, nil
}
