package group

import (
	"context"

	"github.com/google/uuid"
)

type GroupRepo interface {
	InsertGroup(context.Context, *Group) (*Group, error)
	UpdateGroup(context.Context, *Search, *Group) error
	ListGroups(context.Context, *Search) ([]*Group, error)
	SelectGroup(context.Context, *Search) (*Group, error)
	TotalGroup(context.Context, *Search) (*int64, error)
	CheckCategoryIdExisted(context.Context, *Search, uuid.UUID) (*Group, error)
}
