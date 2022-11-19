package api

import (
	"context"

	"github.com/aqaurius6666/apiservice/src/internal/db/category"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/model"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/apiservice/src/pb"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/aqaurius6666/go-utils/utils"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

type CategoryService struct {
	Model model.Server
}

func (s *CategoryService) ListCategories(ctx context.Context, req *pb.CategoriesGetRequest) (*pb.CategoriesGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListCategories))
	defer span.End()
	limit := lib.ParseInt32Val(req.Limit)
	offset := lib.ParseInt32Val(req.Offset)

	total, err := s.Model.TotalCategory(ctx, &category.Search{
		Query: req.Query,
		Category: category.Category{
			Name: utils.SafeStrPtr(lib.StandardizeSpaces(req.Name)),
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	r, err := s.Model.ListCategories(ctx, &category.Search{
		DefaultSearchModel: database.DefaultSearchModel{
			Skip:   int(offset),
			Limit:  int(limit),
			Fields: []string{`"categories"."id"`, `"categories"."name"`, `"categories"."image_url"`},
		},
		Query: req.Query,
		Category: category.Category{
			Name: utils.SafeStrPtr(lib.StandardizeSpaces(req.Name)),
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &pb.CategoriesGetResponse_Data{
		Result:     s.Model.ConvertCategoriesToProtos(r),
		Pagination: lib.Pagination(offset, limit, total),
	}, nil
}

func (s *CategoryService) GetCategory(ctx context.Context, req *pb.CategoryGetRequest) (*pb.CategoryGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetCategory))
	defer span.End()

	r, err := s.Model.GetCategoryById(ctx, req.Id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &pb.CategoryGetResponse_Data{
		Category: s.Model.ConvertCategoryToProto(r),
	}, nil
}
