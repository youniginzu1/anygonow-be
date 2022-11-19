package model

import (
	"context"

	"github.com/aqaurius6666/apiservice/src/internal/db/category"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/apiservice/src/pb"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

var (
	_ CategoryModel = (*ServerModel)(nil)
)

type CategoryModel interface {
	ListCategories(context.Context, *category.Search) ([]*category.Category, error)
	GetCategoryById(context.Context, interface{}) (*category.Category, error)
	TotalCategory(context.Context, *category.Search) (*int64, error)

	ConvertCategoryToProto(*category.Category) *pb.Category
	ConvertCategoriesToProtos([]*category.Category) []*pb.Category

	GetCategoryByName(ctx context.Context, name *string) (*category.Category, error)
	InsertCategory(ctx context.Context, category *category.Category) error
	EditCategory(ctx context.Context, search *category.Search, category *category.Category) error
	DeleteCategory(ctx context.Context, serviceId interface{}) error

	GetCategories(ctx context.Context, search *category.Search) ([]*category.Category, error)
}

func (s *ServerModel) GetCategoryById(ctx context.Context, id interface{}) (*category.Category, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetCategoryById))
	defer span.End()

	uid, err := lib.ToUUID(id)
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}
	r, err := s.Repo.SelectCategory(ctx, &category.Search{
		Category: category.Category{
			BaseModel: database.BaseModel{ID: uid},
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return r, nil
}

func (s *ServerModel) ListCategories(ctx context.Context, search *category.Search) ([]*category.Category, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListCategories))
	defer span.End()

	r, err := s.Repo.ListCategorys(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return r, nil
}

func (s *ServerModel) ConvertCategoriesToProtos(u []*category.Category) []*pb.Category {
	arr := make([]*pb.Category, 0)
	for _, uu := range u {
		arr = append(arr, s.ConvertCategoryToProto(uu))
	}
	return arr
}
func (s *ServerModel) ConvertCategoryToProto(u *category.Category) *pb.Category {
	upb := new(pb.Category)
	if u.ID != uuid.Nil {
		upb.Id = u.ID.String()
	}
	if u.Name != nil {
		upb.Name = *u.Name
	}
	if u.TotalProvider != nil {
		upb.TotalProvider = *u.TotalProvider
	}
	if u.Fee != nil {
		upb.Fee = *u.Fee
	}
	if u.ImageUrl != nil {
		upb.Image = *u.ImageUrl
	}
	return upb
}

func (s *ServerModel) GetCategoryByName(ctx context.Context, name *string) (*category.Category, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetCategoryById))
	defer span.End()

	r, err := s.Repo.SelectCategory(ctx, &category.Search{
		Category: category.Category{
			Name: name,
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return r, nil
}
func (s *ServerModel) GetCategories(ctx context.Context, search *category.Search) ([]*category.Category, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetCategories))
	defer span.End()
	r, err := s.Repo.ListCategoriesAdmin(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return r, nil
}

func (s *ServerModel) InsertCategory(ctx context.Context, category *category.Category) error {

	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetCategoryById))
	defer span.End()
	_, err := s.Repo.InsertCategory(ctx, category)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}

func (s *ServerModel) DeleteCategory(ctx context.Context, serviceId interface{}) error {

	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.DeleteCategory))
	defer span.End()
	sid, err := lib.ToUUID(serviceId)
	if err != nil {
		return xerrors.Errorf("%w", err)
	}
	err = s.Repo.DeleteCategory(ctx, &category.Category{
		BaseModel: database.BaseModel{
			ID: sid,
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}

func (s *ServerModel) EditCategory(ctx context.Context, search *category.Search, category *category.Category) error {

	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.EditCategory))
	defer span.End()

	err := s.Repo.UpdateCategory(ctx, search, category)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}
