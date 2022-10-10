package api

import (
	"context"

	"github.com/aqaurius6666/authservice/src/internal/lib"
	"github.com/aqaurius6666/authservice/src/internal/var/c"
	"github.com/aqaurius6666/authservice/src/pb/authpb"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

func (s *ApiServer) ChangePassword(ctx context.Context, req *authpb.ChangePasswordRequest) (*authpb.ChangePasswordResponse, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ChangePassword))
	defer span.End()
	err := s.Model.UpdateKey(ctx, req.Id, req.PublicKey, req.EncryptedPrivateKey)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		panic(err)
	}
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		panic(err)
	}

	return &authpb.ChangePasswordResponse{}, nil
}
