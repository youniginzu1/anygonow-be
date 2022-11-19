package business

import (
	"context"
)

type BusinessRepo interface {
	SelectBusiness(context.Context, *Search) (*Business, error)
	InsertBusiness(context.Context, *Business) (*Business, error)
	UpdateBusiness(context.Context, *Search, *Business) error
	ListBusinesss(context.Context, *Search) ([]*Business, error)
	ListBusinessOptimize(context.Context, *Search) ([]*Business, error)
	DeleteBusiness(context.Context, *Search) error
	TotalBusiness(context.Context, *Search) (*int64, error)
	ListBusinesssWithRating(context.Context, *Search) ([]*Business, error)
	GetMapIdName(context.Context, *Search) (map[string]*string, error)
	GetTotalZipcodes(context.Context, *Search) (*int64, error)
}
