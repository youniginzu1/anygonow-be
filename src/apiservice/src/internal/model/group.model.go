package model

import (
	"context"

	"github.com/aqaurius6666/apiservice/src/internal/db/group"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/apiservice/src/pb"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

var (
	_ GroupModel = (*ServerModel)(nil)
)

type GroupModel interface {
	InsertGroup(ctx context.Context, group *group.Group) error
	EditGroup(ctx context.Context, search *group.Search, category *group.Group) error
	ListGroups(ctx context.Context, search *group.Search) ([]*group.Group, error)
	SelectGroup(ctx context.Context, search *group.Search) (*group.Group, error)
	TotalGroup(ctx context.Context, search *group.Search) (*int64, error)
	CheckCategoryIdExisted(ctx context.Context, search *group.Search, categoryId uuid.UUID) error
	ConvertGroupToProtos(v []*group.Group) []*pb.Group
	ConvertGroupToProto(g *group.Group) *pb.Group
}

func (s *ServerModel) ConvertGroupToProto(g *group.Group) *pb.Group {
	gpb := &pb.Group{}

	if g.ID != uuid.Nil {
		gpb.Id = g.ID.String()
	}
	if g.Name != nil {
		gpb.Name = *g.Name
	}
	if g.Fee != nil {
		gpb.Fee = *g.Fee
	}
	if g.ServiceName != nil && g.ServiceId != nil {
		for i := range g.ServiceName {
			gpb.ServiceInfo = append(gpb.ServiceInfo, &pb.ServiceGroup{
				ServiceName: g.ServiceName[i],
				ServiceId:   g.ServiceId[i],
			})
		}
	}
	return gpb
}

func (s *ServerModel) ConvertGroupToProtos(v []*group.Group) []*pb.Group {
	arr := make([]*pb.Group, 0)
	for _, vv := range v {
		arr = append(arr, s.ConvertGroupToProto(vv))
	}
	return arr
}

func (s *ServerModel) InsertGroup(ctx context.Context, group *group.Group) error {

	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.InsertGroup))
	defer span.End()

	_, err := s.Repo.InsertGroup(ctx, group)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}

func (s *ServerModel) EditGroup(ctx context.Context, search *group.Search, group *group.Group) error {

	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.EditGroup))
	defer span.End()

	err := s.Repo.UpdateGroup(ctx, search, group)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}

func (s *ServerModel) ListGroups(ctx context.Context, search *group.Search) ([]*group.Group, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListGroups))
	defer span.End()

	gr, err := s.Repo.ListGroups(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return gr, nil
}

func (s *ServerModel) SelectGroup(ctx context.Context, search *group.Search) (*group.Group, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.SelectGroup))
	defer span.End()

	gr, err := s.Repo.SelectGroup(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return gr, nil
}

func (s *ServerModel) TotalGroup(ctx context.Context, search *group.Search) (*int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.TotalGroup))
	defer span.End()

	gr, err := s.Repo.TotalGroup(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return gr, nil
}

func (s *ServerModel) CheckCategoryIdExisted(ctx context.Context, search *group.Search, categoryId uuid.UUID) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CheckCategoryIdExisted))
	defer span.End()

	_, err := s.Repo.CheckCategoryIdExisted(ctx, search, categoryId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}
