package model

import (
	"context"

	"github.com/aqaurius6666/apiservice/src/internal/db/advertise_package"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/apiservice/src/pb"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

var (
	_ AdvertisePackageModel = (*ServerModel)(nil)
)

type AdvertisePackageModel interface {
	InsertPackage(ctx context.Context, adpackage *advertise_package.AdvertisePackage) error
	TotalPackage(ctx context.Context, search *advertise_package.Search) (*int64, error)
	ListPackages(ctx context.Context, search *advertise_package.Search) ([]*advertise_package.AdvertisePackage, error)
	ConvertPackageToProto(p *advertise_package.AdvertisePackage) *pb.AdvertisePackage
	ConvertPackageToProtos(v []*advertise_package.AdvertisePackage) []*pb.AdvertisePackage
	EditAdvertisePackage(ctx context.Context, search *advertise_package.Search, p *advertise_package.AdvertisePackage) error
	DeleteAdvertisePackage(ctx context.Context, id interface{}) error
	CheckCateIdExistedAdvertise(ctx context.Context, search *advertise_package.Search, categoryId uuid.UUID) error

	ListAdvertiseDetails(ctx context.Context, search *advertise_package.Search) ([]*advertise_package.AdvertisePackage, error)
	ConvertAdvertiseDetailToProto(p *advertise_package.AdvertisePackage) *pb.AdvertiseDetail
	ConvertAdvertiseDetailToProtos(v []*advertise_package.AdvertisePackage) []*pb.AdvertiseDetail
	TotalAdvertiseDetail(ctx context.Context, search *advertise_package.Search) (*int64, error)
}

func (s *ServerModel) InsertPackage(ctx context.Context, adpackage *advertise_package.AdvertisePackage) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.InsertPackage))
	defer span.End()

	_, err := s.Repo.InsertAdvertisePackage(ctx, adpackage)

	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}

	return nil
}

func (s *ServerModel) ConvertPackageToProto(p *advertise_package.AdvertisePackage) *pb.AdvertisePackage {
	ppb := &pb.AdvertisePackage{}

	if p.ID != uuid.Nil {
		ppb.Id = p.ID.String()
	}

	if p.Name != nil {
		ppb.Name = *p.Name
	}

	if p.Price != nil {
		ppb.Price = *p.Price
	}

	if p.ServiceName != nil && p.ServiceId != nil {
		for i := range p.ServiceName {
			ppb.ServiceInfo = append(ppb.ServiceInfo, &pb.ServiceGroup{
				ServiceName: p.ServiceName[i],
				ServiceId:   p.ServiceId[i],
			})
		}
	}

	if p.BannerUrl != nil {
		ppb.BannerUrl = *p.BannerUrl
	}

	if p.Description != nil {
		ppb.Description = *p.Description
	}

	return ppb
}

func (s *ServerModel) ConvertPackageToProtos(v []*advertise_package.AdvertisePackage) []*pb.AdvertisePackage {
	arr := make([]*pb.AdvertisePackage, 0)
	for _, vv := range v {
		arr = append(arr, s.ConvertPackageToProto(vv))
	}
	return arr
}

func (s *ServerModel) TotalPackage(ctx context.Context, search *advertise_package.Search) (*int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.TotalPackage))
	defer span.End()

	pk, err := s.Repo.TotalPackage(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return pk, nil
}

func (s *ServerModel) ListPackages(ctx context.Context, search *advertise_package.Search) ([]*advertise_package.AdvertisePackage, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListPackages))
	defer span.End()

	pk, err := s.Repo.ListPackages(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return pk, nil
}

func (s *ServerModel) EditAdvertisePackage(ctx context.Context, search *advertise_package.Search, p *advertise_package.AdvertisePackage) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.EditAdvertisePackage))
	defer span.End()

	err := s.Repo.UpdateAdvertisePackage(ctx, search, p)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}

func (s *ServerModel) DeleteAdvertisePackage(ctx context.Context, id interface{}) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.DeleteAdvertisePackage))
	defer span.End()

	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}

	err = s.Repo.DeleteAdvertisePackage(ctx, &advertise_package.AdvertisePackage{
		BaseModel: database.BaseModel{ID: uid},
	})

	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}

	return nil
}

func (s *ServerModel) CheckCateIdExistedAdvertise(ctx context.Context, search *advertise_package.Search, categoryId uuid.UUID) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CheckCateIdExistedAdvertise))
	defer span.End()

	_, err := s.Repo.CheckCateIdExistedAdvertise(ctx, search, categoryId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}

func (s *ServerModel) ListAdvertiseDetails(ctx context.Context, search *advertise_package.Search) ([]*advertise_package.AdvertisePackage, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListAdvertiseDetails))
	defer span.End()

	pk, err := s.Repo.ListAdvertiseDetails(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return pk, nil
}

func (s *ServerModel) ConvertAdvertiseDetailToProto(p *advertise_package.AdvertisePackage) *pb.AdvertiseDetail {
	ppb := &pb.AdvertiseDetail{}

	if p.ID != uuid.Nil {
		ppb.Id = p.ID.String()
	}

	if p.Name != nil {
		ppb.Name = *p.Name
	}

	if p.Price != nil {
		ppb.Price = *p.Price
	}

	if p.ServiceName != nil && p.ServiceId != nil {
		for i := range p.ServiceName {
			ppb.ServiceInfo = append(ppb.ServiceInfo, &pb.ServiceGroup{
				ServiceName: p.ServiceName[i],
				ServiceId:   p.ServiceId[i],
			})
		}
	}

	if p.BannerUrl != nil {
		ppb.BannerUrl = *p.BannerUrl
	}

	if p.Description != nil {
		ppb.Description = *p.Description
	}

	if p.Zipcodes != nil {
		for i := range p.Zipcodes {
			ppb.Zipcodes = append(ppb.Zipcodes, p.Zipcodes[i])
		}
	}

	return ppb
}

func (s *ServerModel) ConvertAdvertiseDetailToProtos(v []*advertise_package.AdvertisePackage) []*pb.AdvertiseDetail {
	arr := make([]*pb.AdvertiseDetail, 0)
	for _, vv := range v {
		arr = append(arr, s.ConvertAdvertiseDetailToProto(vv))
	}
	return arr
}

func (s *ServerModel) TotalAdvertiseDetail(ctx context.Context, search *advertise_package.Search) (*int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.TotalAdvertiseDetail))
	defer span.End()

	pk, err := s.Repo.TotalAdvertiseDetail(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return pk, nil
}
