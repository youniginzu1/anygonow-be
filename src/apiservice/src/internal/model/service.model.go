package model

import (
	"context"

	"github.com/aqaurius6666/apiservice/src/internal/db/business"
	"github.com/aqaurius6666/apiservice/src/internal/db/category"
	"github.com/aqaurius6666/apiservice/src/internal/db/service"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/apiservice/src/pb"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/aqaurius6666/go-utils/utils"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

var (
	_ ServiceModel = (*ServerModel)(nil)
)

type ServiceModel interface {
	GetServicesByBusiness(ctx context.Context, id interface{}) ([]*service.Service, error)
	GetServiceById(ctx context.Context, id interface{}) (*service.Service, error)
	UpdateServices(ctx context.Context, id interface{}, categoryIds []uuid.UUID) ([]*service.Service, error)

	ConvertServiceToProto(*service.Service) *pb.Service
	ConvertServiceToProtos([]*service.Service) []*pb.Service
	TotalCategory(ctx context.Context, search *category.Search) (*int64, error)
}

func (s *ServerModel) UpdateServices(ctx context.Context, id interface{}, categoryIds []uuid.UUID) ([]*service.Service, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UpdateServices))
	defer span.End()

	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	countCate, err := s.TotalCategory(ctx, &category.Search{
		CategoryIds: categoryIds,
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	if int(*countCate) != len(categoryIds) {
		err = xerrors.Errorf("%w", category.ErrNotFound)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	if err = s.Repo.UpdateService(ctx, &service.Search{
		Service: service.Service{
			BusinessId: uid,
		},
	}, &service.Service{
		Status: utils.Int32Ptr(int32(c.SERVICE_STATUS_SERVICE_INACTIVE)),
	}); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	serviceIds := make([]uuid.UUID, 0)
	for _, cate := range categoryIds {
		s, err := s.Repo.FirstOrInsertService(ctx, &service.Search{
			Service: service.Service{
				BusinessId: uid,
				CategoryId: cate,
			},
		}, &service.Service{
			BusinessId: uid,
			CategoryId: cate,
			Status:     utils.Int32Ptr(int32(c.SERVICE_STATUS_SERVICE_INACTIVE)),
		})
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err, ctx)
			return nil, err
		}
		serviceIds = append(serviceIds, s.ID)

	}
	if err = s.Repo.UpdateService(ctx, &service.Search{
		Service: service.Service{
			BusinessId: uid,
		},
		CategoryIds: categoryIds,
	}, &service.Service{
		Status: utils.Int32Ptr(int32(c.SERVICE_STATUS_SERVICE_ACTIVE)),
	}); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	if err = s.UpdateBusiness(ctx, uid, &business.Business{
		Services: serviceIds,
	}); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return nil, nil
}

func (s *ServerModel) GetServiceById(ctx context.Context, id interface{}) (*service.Service, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetServiceById))
	defer span.End()

	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	ser, err := s.Repo.SelectService(ctx, &service.Search{
		Service: service.Service{
			BaseModel: database.BaseModel{ID: uid},
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return ser, nil
}

func (s *ServerModel) GetServicesByBusiness(ctx context.Context, id interface{}) ([]*service.Service, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetServicesByBusiness))
	defer span.End()

	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	services, err := s.Repo.ListServices(ctx, &service.Search{
		Service: service.Service{
			BusinessId: uid,
			Status:     utils.Int32Ptr(int32(c.SERVICE_STATUS_SERVICE_ACTIVE)),
		},
		DefaultSearchModel: database.DefaultSearchModel{
			Fields: []string{
				`"categories"."name" as "name"`,
				`"services"."id"`,
				`"categories"."image_url" as "logo_url"`,
				`count("orders"."id") as "number_order"`,
				`"categories"."id" as "category_id"`,
			},
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return services, nil
}

func (s *ServerModel) ConvertServiceToProto(svc *service.Service) *pb.Service {
	upb := &pb.Service{}
	if svc.ID != uuid.Nil {
		upb.Id = svc.ID.String()
	}
	if svc.Name != nil {
		upb.Name = utils.StrVal(svc.Name)
	}
	if svc.LogoUrl != nil {
		upb.Image = utils.StrVal(svc.LogoUrl)
	}
	if svc.BusinessId != uuid.Nil {
		upb.BusinessId = svc.BusinessId.String()
	}
	if svc.Status != nil {
		upb.Status = c.SERVICE_STATUS(utils.Int32Val(svc.Status))
	}
	if svc.CategoryId != uuid.Nil {
		upb.CategoryId = svc.CategoryId.String()
	}
	if svc.Category != nil && svc.Category.Name != nil {
		upb.CategoryName = *svc.Category.Name
	}
	if svc.NumberOrder != nil {
		upb.NumberOrder = *svc.NumberOrder
	}
	return upb
}

func (s *ServerModel) ConvertServiceToProtos(svc []*service.Service) []*pb.Service {
	upb := make([]*pb.Service, 0)
	for _, sv := range svc {
		upb = append(upb, s.ConvertServiceToProto(sv))
	}
	return upb
}

func (s *ServerModel) TotalCategory(ctx context.Context, search *category.Search) (*int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.TotalCategory))
	defer span.End()

	r, err := s.Repo.TotalCategory(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return r, nil
}
