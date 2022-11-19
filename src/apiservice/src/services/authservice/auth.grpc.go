package authservice

import (
	"context"
	"time"

	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/apiservice/src/internal/var/e"
	"github.com/aqaurius6666/apiservice/src/pb/authpb"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type AuthServiceAddr string

var (
	_       Service = ServiceGRPC{}
	timeout         = 5 * time.Second
)

type ServiceGRPC struct {
	Ctx    context.Context
	Client authpb.AuthServiceClient
}

func ConnectClient(ctx context.Context, addr AuthServiceAddr) (authpb.AuthServiceClient, error) {
	nctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	conn, err := grpc.DialContext(nctx, string(addr), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithChainUnaryInterceptor(
		otelgrpc.UnaryClientInterceptor(),
	))
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}
	return authpb.NewAuthServiceClient(conn), nil
}
func (s ServiceGRPC) ResendOTP(ctx context.Context, id string) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ResendOTP))
	defer span.End()
	_, err := s.Client.ResendOTP(ctx, &authpb.ResendOTPRequest{
		Id: id,
	})
	if err != nil {
		stt, ok := status.FromError(err)
		if ok {
			err = xerrors.New(stt.Message())
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}
func (s ServiceGRPC) CheckCredential(ctx context.Context, id string) (bool, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CheckCredential))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	res, err := s.Client.CheckCredential(ctx, &authpb.CheckCredentialRequest{
		Identifier: id,
	})
	if err != nil {
		stt, ok := status.FromError(err)
		if ok {
			err = xerrors.New(stt.Message())
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return false, err
	}
	return res.Existed, nil
}

func (s ServiceGRPC) DeleteUser(ctx context.Context, id string) (string, c.ROLE, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.DeleteUser))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	res, err := s.Client.DeleteUser(ctx, &authpb.DeleteUserRequest{
		Id: id,
	})
	if err != nil {
		stt, ok := status.FromError(err)
		if ok {
			err = xerrors.New(stt.Message())
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return "", 0, err
	}
	return res.Id, res.Role, nil
}
func (s ServiceGRPC) BanUser(ctx context.Context, id string, stt bool) (string, c.ROLE, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.BanUser))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	res, err := s.Client.BanUser(ctx, &authpb.BanUserRequest{
		Id:     id,
		Status: stt,
	})
	if err != nil {
		stt, ok := status.FromError(err)
		if ok {
			err = xerrors.New(stt.Message())
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return "", 0, err
	}
	return res.Id, res.Role, nil
}

func (s ServiceGRPC) ForgotPassword(ctx context.Context, mail string) (string, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ForgotPassword))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	res, err := s.Client.ForgotPassword(ctx, &authpb.ForgotPasswordRequest{
		Mail: mail,
	})
	if err != nil {
		stt, ok := status.FromError(err)
		if ok {
			err = xerrors.New(stt.Message())
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)

		return "", err
	}
	return res.OtpId, nil
}

func (s ServiceGRPC) ChangeMail(ctx context.Context, id string, mail string) (string, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ChangeMail))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	res, err := s.Client.ChangeMail(ctx, &authpb.ChangeMailRequest{
		Mail: mail,
		Id:   id,
	})
	if err != nil {
		stt, ok := status.FromError(err)
		if ok {
			err = xerrors.New(stt.Message())
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return "", err
	}
	return res.OtpId, nil
}

func (s ServiceGRPC) ChangeMailOTP(ctx context.Context, otpId string) (string, string, c.ROLE, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ChangeMailOTP))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	res, err := s.Client.ChangeMailOTP(ctx, &authpb.ChangeMailOTPRequest{
		OtpId: otpId,
	})
	if err != nil {
		stt, ok := status.FromError(err)
		if ok {
			err = xerrors.New(stt.Message())
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return "", "", 0, err
	}
	return res.Id, res.Mail, res.Role, nil
}

func (s ServiceGRPC) ChangeMailAndPassOTP(ctx context.Context, otpId string) (string, string, c.ROLE, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ChangeMailAndPass))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	res, err := s.Client.ChangeMailAndPassOTP(ctx, &authpb.ChangeMailAndPassOTPRequest{
		OtpId: otpId,
	})
	if err != nil {
		stt, ok := status.FromError(err)
		if ok {
			err = xerrors.New(stt.Message())
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return "", "", 0, err
	}
	return res.Id, res.Mail, res.Role, nil
}

func (s ServiceGRPC) ForgotPasswordOTP(ctx context.Context, otpId, pub, enc string) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ForgotPasswordOTP))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	_, err := s.Client.ForgotPasswordOTP(ctx, &authpb.ForgotPasswordOTPRequest{
		OtpId:               otpId,
		EncryptedPrivateKey: enc,
		PublicKey:           pub,
	})
	if err != nil {
		stt, ok := status.FromError(err)
		if ok {
			err = xerrors.New(stt.Message())
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}

func (s ServiceGRPC) RegisterOTP(ctx context.Context, otpId string) (id, mail, phone string, role c.ROLE, refCode string, err error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.RegisterOTP))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	res, err := s.Client.RegisterOTP(ctx, &authpb.RegisterOTPRequest{
		OtpId: otpId,
	})
	if err != nil {
		stt, ok := status.FromError(err)
		if ok {
			err = xerrors.New(stt.Message())
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return
	}
	return res.Id, res.Mail, res.Phone, res.Role, res.RefCode, nil
}

func (s ServiceGRPC) VerifyOTP(ctx context.Context, otpId, otp string) (typ c.OTP_TYPE, err error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.VerifyOTP))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	res, err := s.Client.VerifyOTP(ctx, &authpb.VerifyOTPRequest{
		OtpId: otpId,
		Otp:   otp,
	})
	if err != nil {
		stt, ok := status.FromError(err)
		if ok {
			err = xerrors.New(stt.Message())
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return
	}
	if !res.Ok {
		err = xerrors.Errorf("%w", e.ErrOTPFail)
		lib.RecordError(span, err, ctx)
		return 0, err
	}
	return res.Type, nil
}

func (s ServiceGRPC) CheckAuth(ctx context.Context, header, body []byte, method string) (string, c.ROLE, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CheckAuth))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	res, err := s.Client.CheckAuth(ctx, &authpb.CheckAuthRequest{
		Header: header,
		Body:   body,
		Method: method,
	})
	if err != nil {
		stt, ok := status.FromError(err)
		if ok {
			err = xerrors.New(stt.Message())
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return "", 0, err
	}
	return res.Id, res.Role, nil
}

func (s ServiceGRPC) Register(ctx context.Context, mail, phone, pub, enpriv string, role c.ROLE, refCode string) (otpId string, err error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.Register))
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	res, err := s.Client.Register(ctx, &authpb.RegisterRequest{
		Username:            mail,
		Mail:                mail,
		Phone:               phone,
		PublicKey:           pub,
		EncryptedPrivateKey: enpriv,
		Role:                role,
		RefCode:             refCode,
	})

	if err != nil {
		stt, ok := status.FromError(err)
		if ok {
			err = xerrors.New(stt.Message())
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return
	}
	return res.OtpId, nil
}

func (s ServiceGRPC) GetCredential(ctx context.Context, data string) (id, pub, priv string, isDef bool, err error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetCredential))
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	res, err := s.Client.GetCredential(ctx, &authpb.GetCredentialRequest{
		Identifier: data,
	})
	if err != nil {
		stt, ok := status.FromError(err)
		if ok {
			err = xerrors.New(stt.Message())
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return
	}
	return res.Credential.Id, res.Credential.PublicKey, res.Credential.EncryptedPrivateKey, res.Credential.IsDefaultPassword, nil
}

func (s ServiceGRPC) ChangePassword(ctx context.Context, id, pub, enc string) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ChangePassword))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := s.Client.ChangePassword(ctx, &authpb.ChangePasswordRequest{
		Id:                  id,
		PublicKey:           pub,
		EncryptedPrivateKey: enc,
	})
	if err != nil {
		stt, ok := status.FromError(err)
		if ok {
			err = xerrors.New(stt.Message())
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}

func (s ServiceGRPC) ChangeMailAndPass(ctx context.Context, userId, mail, pub, enc string) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ChangePassword))
	defer span.End()

	_, err := s.Client.ChangeMailAndPass(ctx, &authpb.ChangeMailAndPassRequest{
		Mail:                mail,
		EncryptedPrivateKey: enc,
		UserId:              userId,
		PublicKey:           pub,
	})
	if err != nil {
		stt, ok := status.FromError(err)
		if ok {
			err = xerrors.New(stt.Message())
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}
