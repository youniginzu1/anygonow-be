package api

import (
	"context"
	"strings"

	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/lib/validate"
	"github.com/aqaurius6666/apiservice/src/internal/model"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/apiservice/src/internal/var/e"
	"github.com/aqaurius6666/apiservice/src/pb"
	"github.com/aqaurius6666/go-utils/utils"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

type AuthService struct {
	Model model.Server
}

func (s AuthService) CheckCredential(ctx context.Context, req *pb.AuthCheckGetRequest) (*pb.AuthCheckGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CheckCredential))
	defer span.End()
	if f, ok := validate.RequiredFields(req, "Identifier"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	existed, err := s.Model.CheckCredential(ctx, req.Identifier)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &pb.AuthCheckGetResponse_Data{
		Existed: existed,
	}, nil
}

func (s AuthService) GetCredential(ctx context.Context, req *pb.AuthCredentialRequest) (*pb.AuthCredentialResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetCredential))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "Identifier"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	id, pub, priv, isDef, err := s.Model.GetCredential(ctx, req.Identifier)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	changeMail := "@anygonow.com"
	shouldChangeMail := strings.Contains(req.Identifier, changeMail)

	return &pb.AuthCredentialResponse_Data{
		Id:                  id,
		PublicKey:           pub,
		EncryptedPrivateKey: priv,
		ShouldChangeMail:    shouldChangeMail,
		IsDefaultPassword:   isDef,
	}, nil
}

func (s AuthService) ChangeMail(ctx context.Context, req *pb.AuthMailPostRequest) (*pb.AuthMailPostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ChangeMail))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId", "Mail"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	otpId, err := s.Model.ChangeMail(ctx, req.XUserId, req.Mail)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &pb.AuthMailPostResponse_Data{
		OtpId: otpId,
	}, nil
}

func (s AuthService) VerifyOTP(ctx context.Context, req *pb.AuthOTPPostRequest) (*pb.AuthOTPPostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.VerifyOTP))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "OtpId", "Otp"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	err := s.Model.VerifyOTP(ctx, req.OtpId, req.Otp)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.AuthOTPPostResponse_Data{}, nil
}

func (s AuthService) ResendOTP(ctx context.Context, req *pb.AuthResendOTPPostRequest) (*pb.AuthResendOTPPostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ResendOTP))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "OtpId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	err := s.Model.ResendOtp(ctx, req.OtpId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.AuthResendOTPPostResponse_Data{}, nil
}

func (s AuthService) Ping(ctx context.Context, req *pb.AuthPingRequest) (*pb.AuthPingResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.Ping))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err

	}
	switch req.XRole {
	case c.ROLE_HANDYMAN:
		process, err := s.Model.GetRegistrationProccess(ctx, req.XUserId)
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err, ctx)
			return nil, err
		}
		bus, err := s.Model.GetBusinessById(ctx, req.XUserId)
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err, ctx)
			return nil, err
		}
		return &pb.AuthPingResponse_Data{
			Id:       req.XUserId,
			Role:     req.XRole,
			Process:  process,
			Image:    utils.StrVal(bus.LogoUrl),
			Business: s.Model.ConvertBusinessToProto(bus),
		}, nil
	case c.ROLE_CUSTOMER:
		process, err := s.Model.GetRegistrationProccessForUser(ctx, req.XUserId)
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err, ctx)
			return nil, err
		}
		usr, err := s.Model.GetUserById(ctx, req.XUserId)
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err, ctx)
			return nil, err
		}
		return &pb.AuthPingResponse_Data{
			Id:      req.XUserId,
			Role:    req.XRole,
			Process: process,
			Image:   utils.StrVal(usr.AvatarUrl),
			User:    s.Model.ConvertUserToProto(usr),
		}, nil
	case c.ROLE_ADMIN:
		return &pb.AuthPingResponse_Data{
			Id:      req.XUserId,
			Role:    req.XRole,
			Process: c.REGISTRATION_PROCESS_DONE,
		}, nil
	default:
		return nil, nil
	}

}

func (s AuthService) ChangePassword(ctx context.Context, req *pb.AuthPasswordPostRequest) (*pb.AuthPasswordPostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ChangePassword))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId", "PublicKey", "EncryptedPrivateKey"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	if err := s.Model.ChangePassword(ctx, req.XUserId, req.PublicKey, req.EncryptedPrivateKey); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &pb.AuthPasswordPostResponse_Data{}, nil
}

func (s AuthService) ForgotPassword(ctx context.Context, req *pb.AuthForgotPostRequest) (*pb.AuthForgotPostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ForgotPassword))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "Mail"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	id, _, _, _, err := s.Model.GetCredential(ctx, req.Mail)
	if err != nil {
		return &pb.AuthForgotPostResponse_Data{}, nil
	}

	otpId, err := s.Model.ForgotPassword(ctx, req.Mail)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.AuthForgotPostResponse_Data{
		OtpId: otpId,
		Id:    id,
	}, nil
}

func (s AuthService) ResetPassword(ctx context.Context, req *pb.AuthForgotResetPostRequest) (*pb.AuthForgotResetPostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ResetPassword))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "OtpId", "PublicKey", "EncryptedPrivateKey"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	if err := s.Model.ForgotPasswordOTP(ctx, req.OtpId, req.PublicKey, req.EncryptedPrivateKey); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.AuthForgotResetPostResponse_Data{}, nil
}

func (s AuthService) ChangeMailAndPass(ctx context.Context, req *pb.AuthChangeMailAndPassPostRequest) (*pb.AuthChangeMailAndPassPostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ChangeMailAndPass))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId", "Mail", "PublicKey", "EncryptedPrivateKey"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	if err := s.Model.ChangePassAndMail(ctx, req.XUserId, req.Mail, req.PublicKey, req.EncryptedPrivateKey); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.AuthChangeMailAndPassPostResponse_Data{}, nil
}
