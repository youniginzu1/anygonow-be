package advertise_package

import (
	"context"

	"github.com/google/uuid"
)

type AdvertisePackageRepo interface {
	SelectAdvertisePackage(context.Context, *Search) (*AdvertisePackage, error)
	InsertAdvertisePackage(context.Context, *AdvertisePackage) (*AdvertisePackage, error)
	UpdateAdvertisePackage(context.Context, *Search, *AdvertisePackage) error
	TotalPackage(context.Context, *Search) (*int64, error)
	ListPackages(context.Context, *Search) ([]*AdvertisePackage, error)
	DeleteAdvertisePackage(context.Context, *AdvertisePackage) error
	CheckCateIdExistedAdvertise(context.Context, *Search, uuid.UUID) (*AdvertisePackage, error)
	ListAdvertiseDetails(context.Context, *Search) ([]*AdvertisePackage, error)
	TotalAdvertiseDetail(context.Context, *Search) (*int64, error)
}
