package model

import (
	"context"

	"github.com/aqaurius6666/authservice/src/internal/db/role"
	"github.com/aqaurius6666/authservice/src/internal/lib"
	"github.com/aqaurius6666/authservice/src/internal/var/c"
	"github.com/aqaurius6666/go-utils/utils"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

type RoleModel interface {
	GetRoleByCode(context.Context, int) (*role.Role, error)
}

func (s *ServerModel) GetRoleByCode(ctx context.Context, i int) (*role.Role, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetRoleByCode))
	defer span.End()

	r, err := s.Repo.SelectRole(ctx, &role.Search{
		Role: role.Role{
			Code: utils.IntPtr(i),
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	return r, nil
}
