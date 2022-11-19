package api

import (
	"context"

	"github.com/aqaurius6666/apiservice/src/internal/db/user"
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

type UserService struct {
	Model model.Server
}

func (s UserService) CreateUser(ctx context.Context, req *pb.UserPostRequest) (*pb.UserPostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CreateUser))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "Mail", "Phone", "PublicKey", "EncryptedPrivateKey"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	if err := req.Validate(); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	otpId, err := s.Model.Register(ctx, req.Mail, req.Phone, req.PublicKey, req.EncryptedPrivateKey, c.ROLE_CUSTOMER, "")
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &pb.UserPostResponse_Data{
		OtpId: otpId,
		Mail:  req.Mail,
		Phone: req.Phone,
	}, nil
}

func (s UserService) GetById(ctx context.Context, req *pb.UserGetRequest) (*pb.UserGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetById))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "Id"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	usr, err := s.Model.GetUserById(ctx, req.Id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &pb.UserGetResponse_Data{
		User: s.Model.ConvertUserToProto(usr),
	}, nil
}

func (s UserService) Update(ctx context.Context, req *pb.UserPutRequest) (*pb.UserPutResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.Update))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "Id", "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	if req.Id != req.XUserId {
		err := xerrors.Errorf("%w", e.ErrNoPermission)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	if err := req.ValidateAll(); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	b := &user.User{
		AvatarUrl: utils.SafeStrPtr(req.Url),
		FirstName: utils.StrPtr(req.FirstName),
		LastName:  utils.StrPtr(req.LastName),
		Phone:     utils.StrPtr(req.Phone),
	}
	if err := validate.Validate(b); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	if err := s.Model.UpdateUser(ctx, req.Id, b); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.UserPutResponse_Data{
		User: s.Model.ConvertUserToProto(b),
	}, nil
}

func (s UserService) StateGet(ctx context.Context, req *pb.UserStateGetRequest) (*pb.UserStateGetResponse_Data, error) {
	_, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.StateGet))
	defer span.End()

	return &pb.UserStateGetResponse_Data{
		Results: lib.ListState,
	}, nil
}
