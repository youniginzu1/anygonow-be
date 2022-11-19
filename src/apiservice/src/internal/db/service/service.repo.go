package service

import "context"

type ServiceRepo interface {
	SelectService(context.Context, *Search) (*Service, error)
	InsertService(context.Context, *Service) (*Service, error)
	ListServices(context.Context, *Search) ([]*Service, error)
	UpdateService(context.Context, *Search, *Service) error
	FirstOrInsertService(context.Context, *Search, *Service) (*Service, error)
}
