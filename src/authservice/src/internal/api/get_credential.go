package api

import (
	"context"

	"github.com/aqaurius6666/authservice/src/internal/db/user"
	"github.com/aqaurius6666/authservice/src/internal/lib"
	"github.com/aqaurius6666/authservice/src/internal/lib/validate"
	"github.com/aqaurius6666/authservice/src/internal/var/c"
	"github.com/aqaurius6666/authservice/src/pb/authpb"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

func (s *ApiServer) GetCredential(ctx context.Context, req *authpb.GetCredentialRequest) (*authpb.GetCredentialResponse, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetCredential))
	defer span.End()
	var (
		usr *user.User
		err error
	)
	identifier := req.Identifier
	if validate.IsMail(identifier) {
		usr, err = s.Model.GetUserByMail(ctx, req.Identifier)
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err)
			panic(err)
		}
	} else if validate.IsPhone(identifier) {
		usr, err = s.Model.GetUserByPhone(ctx, req.Identifier)
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err)
			panic(err)
		}
	} else {
		err = xerrors.Errorf("%w", user.ErrNotFound)
		lib.RecordError(span, err)
		panic(err)
	}

	return &authpb.GetCredentialResponse{
		Credential: &authpb.Credential{
			Id:                  usr.ID.String(),
			PublicKey:           *usr.PublicKey,
			EncryptedPrivateKey: *usr.EncryptedPrivateKey,
			IsDefaultPassword:   *usr.IsDefaultPassword,
		},
	}, nil
}
