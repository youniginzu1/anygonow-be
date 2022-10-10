package role

import "context"

type RoleRepo interface {
	SelectRole(context.Context, *Search) (*Role, error)
	InsertRole(context.Context, *Role) (*Role, error)
}
