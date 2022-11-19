package model

import (
	"context"

	"github.com/aqaurius6666/apiservice/src/internal/db/advertise_order"
	"github.com/aqaurius6666/apiservice/src/internal/db/advertise_transaction"
	"github.com/aqaurius6666/apiservice/src/internal/db/service"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/apiservice/src/pb"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

var (
	_ AdvertiseOrderModel = (*ServerModel)(nil)
)

type AdvertiseOrderModel interface {
	TotalAdvertiseOrder(ctx context.Context, search *advertise_order.Search) (*int64, error)
	ListAdvertiseOrders(ctx context.Context, search *advertise_order.Search) ([]*advertise_order.AdvertiseOrder, error)
	ConvertAdvertiseOrderToProto(o *advertise_order.AdvertiseOrder) *pb.AdvertiseOrder
	ConvertAdvertiseOrderToProtos(v []*advertise_order.AdvertiseOrder) []*pb.AdvertiseOrder
	CreateAdvertiseOrder(ctx context.Context, transaction *advertise_transaction.AdvertiseTransaction, order *advertise_order.AdvertiseOrder) error
	GetServiceForPayment(ctx context.Context, search *service.Search) (*service.Service, error)
	GetTotalOrderForBuyValidate(ctx context.Context, search *advertise_order.Search) (*int64, error)
	TotalFeeAdvertise(ctx context.Context, search *advertise_order.Search) (*float64, error)
}

func (s *ServerModel) TotalAdvertiseOrder(ctx context.Context, search *advertise_order.Search) (*int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.TotalAdvertiseOrder))
	defer span.End()

	pk, err := s.Repo.TotalAdvertiseOrder(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return pk, nil
}

func (s *ServerModel) ListAdvertiseOrders(ctx context.Context, search *advertise_order.Search) ([]*advertise_order.AdvertiseOrder, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListAdvertiseOrders))
	defer span.End()

	pk, err := s.Repo.ListAdvertiseOrders(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return pk, nil
}

func (s *ServerModel) ConvertAdvertiseOrderToProto(o *advertise_order.AdvertiseOrder) *pb.AdvertiseOrder {
	opb := &pb.AdvertiseOrder{}

	if o.ID != uuid.Nil {
		opb.Id = o.ID.String()
	}
	if o.AdvertisePackageId != uuid.Nil {
		opb.AdvertisePackageId = o.AdvertisePackageId.String()
	}
	if o.StartDate != nil {
		opb.StartDate = *o.StartDate
	}
	if o.EndDate != nil {
		opb.EndDate = *o.EndDate
	}
	if o.Name != nil {
		opb.Name = *o.Name
	}
	if o.Price != nil {
		opb.Price = *o.Price
	}
	if o.BannerUrl != nil {
		opb.BannerUrl = *o.BannerUrl
	}
	if o.Description != nil {
		opb.Description = *o.Description
	}
	if o.Zipcode != nil {
		opb.Zipcode = *o.Zipcode
	}
	if o.CategoryName != nil {
		opb.CategoryName = *o.CategoryName
	}
	return opb
}

func (s *ServerModel) ConvertAdvertiseOrderToProtos(v []*advertise_order.AdvertiseOrder) []*pb.AdvertiseOrder {
	arr := make([]*pb.AdvertiseOrder, 0)
	for _, vv := range v {
		arr = append(arr, s.ConvertAdvertiseOrderToProto(vv))
	}
	return arr
}

func (s *ServerModel) CreateAdvertiseOrder(ctx context.Context, transaction *advertise_transaction.AdvertiseTransaction, order *advertise_order.AdvertiseOrder) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CreateAdvertiseOrder))
	defer span.End()

	at, err := s.Repo.InsertAdvertiseTransaction(ctx, transaction)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}

	order.AdvertiseTransactionId = at.ID

	_, err = s.Repo.InsertAdvertiseOrder(ctx, order)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}

	return nil
}

func (s *ServerModel) GetServiceForPayment(ctx context.Context, search *service.Search) (*service.Service, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetServiceForPayment))
	defer span.End()

	ser, err := s.Repo.SelectService(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return ser, nil

}

func (s *ServerModel) GetTotalOrderForBuyValidate(ctx context.Context, search *advertise_order.Search) (*int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetTotalOrderForBuyValidate))
	defer span.End()

	total, err := s.Repo.GetTotalOrderForBuyValidate(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return total, nil
}

func (s *ServerModel) TotalFeeAdvertise(ctx context.Context, search *advertise_order.Search) (*float64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetTotalOrderForBuyValidate))
	defer span.End()

	fee, err := s.Repo.TotalFeeAdvertise(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return fee, nil
}
