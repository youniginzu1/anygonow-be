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

func (s *ApiServer) DeleteUser(ctx context.Context, req *authpb.DeleteUserRequest) (*authpb.DeleteUserResponse, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.DeleteUser))
	defer span.End()
	usr, err := s.Model.GetUserById(ctx, req.Id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		panic(err)
	}
	if *usr.IsActive {
		err = xerrors.Errorf("%w", e.ErrDeleteActiveUser)
		lib.RecordError(span, err)
		panic(err)
	}
	if err := s.Model.DeleteUser(ctx, req.Id); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		panic(err)
	}

	return &authpb.DeleteUserResponse{
		Id:   req.Id,
		Role: c.ROLE(utils.IntVal(usr.Role.Code)),
	}, nil
}
