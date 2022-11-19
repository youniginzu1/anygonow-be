package category

import "context"

type CategoryRepo interface {
	SelectCategory(context.Context, *Search) (*Category, error)
	TotalCategory(context.Context, *Search) (*int64, error)
	InsertCategory(context.Context, *Category) (*Category, error)
	UpdateCategory(context.Context, *Search, *Category) error
	ListCategorys(context.Context, *Search) ([]*Category, error)
	DeleteCategory(context.Context, *Category) error
	ListCategoriesAdmin(context.Context, *Search) ([]*Category, error)
}
