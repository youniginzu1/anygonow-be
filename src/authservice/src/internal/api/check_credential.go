package api

import (
	"context"

	"github.com/aqaurius6666/authservice/src/internal/db/user"
	"github.com/aqaurius6666/authservice/src/internal/lib"
	"github.com/aqaurius6666/authservice/src/internal/var/c"
	"github.com/aqaurius6666/authservice/src/pb/authpb"
	"go.opentelemetry.io/otel"
)

func (s *ApiServer) CheckCredential(ctx context.Context, req *authpb.CheckCredentialRequest) (*authpb.CheckCredentialResponse, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CheckCredential))
	defer span.End()
	if err := s.Model.CheckUserExisted(ctx, &user.User{
		Mail:  &req.Identifier,
		Phone: &req.Identifier,
	}); err != nil {
		return &authpb.CheckCredentialResponse{
			Existed: true,
		}, nil
	}
	return &authpb.CheckCredentialResponse{
		Existed: false,
	}, nil
}
