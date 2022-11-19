package model

import (
	"context"
	"time"

	"github.com/aqaurius6666/apiservice/src/internal/db/order"
	"github.com/aqaurius6666/apiservice/src/internal/db/service"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/apiservice/src/internal/var/e"
	"github.com/aqaurius6666/apiservice/src/pb"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/aqaurius6666/go-utils/utils"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

var THREE_DAYS = 3 * 24 * time.Hour

type OrderModel interface {
	CreateOrder(ctx context.Context, uid, sid interface{}, zipcode *string, phone *string, message *string) (*order.Order, error)
	CreateOrderV2(ctx context.Context, uid, bid, cid interface{}, zipcode *string, phone *string, message *string, customerName *string) (*order.Order, error)
	GetCurrentOrderCount(ctx context.Context, uid interface{}, zipcode *string) (*int64, error)
	ListOrders(context.Context, *order.Search) ([]*order.Order, error)
	TotalOrders(context.Context, *order.Search) (*int64, error)
	GetOrderById(context.Context, interface{}) (*order.Order, error)
	UpdateOrderStatusById(ctx context.Context, orderId interface{}, status c.ORDER_STATUS) error
	UpdateOrderById(ctx context.Context, orderId interface{}, value *order.Order) error
	UpdateOrdersStatusByUser(ctx context.Context, userId interface{}, status c.ORDER_STATUS) error
	ListBusinessesAlreadyOrdered(context.Context, *order.Search) ([]*order.Order, error)
	CloseConversation(ctx context.Context, orderId uuid.UUID) error

	ConvertOrderToProtos(u []*order.Order) []*pb.Order
	ConvertOrderToProto(u *order.Order) *pb.Order

	ConvertProjectToProtos(u []*order.Order) []*pb.Project
	ConvertProjectToProto(u *order.Order) *pb.Project

	ListProjects(ctx context.Context, search *order.Search) ([]*order.Order, error)
	TotalProjects(ctx context.Context, search *order.Search) (*int64, error)
	CancelProject(ctx context.Context, search *order.Search) error

	CheckPermissionUpdateOrder(idRequest string, idFromOrder uuid.UUID) bool
}

func (s *ServerModel) CloseConversation(ctx context.Context, orderId uuid.UUID) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CloseConversation))
	defer span.End()
	err := s.Chat.CloseConversation(ctx, orderId.String())
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return err
	}
	return nil
}

func (s *ServerModel) UpdateOrderById(ctx context.Context, orderId interface{}, value *order.Order) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UpdateOrderById))
	defer span.End()

	id, err := lib.ToUUID(orderId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	err = s.Repo.UpdateOrder(ctx, &order.Search{
		Order: order.Order{BaseModel: database.BaseModel{ID: id}},
	}, value)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}

func (s *ServerModel) ConvertProjectToProto(u *order.Order) *pb.Project {
	upb := new(pb.Project)
	if u.ServiceName != nil {
		upb.ServiceName = *u.ServiceName
	}
	if u.CustomerZipcode != nil {
		upb.Zipcode = *u.CustomerZipcode
	}
	if u.NumberOrders != nil {
		upb.Total = *u.NumberOrders
	}
	if u.ServiceAvatar != nil {
		upb.Image = *u.ServiceAvatar
	}
	if u.ServiceId != uuid.Nil {
		upb.ServiceId = u.ServiceId.String()
	}
	return upb
}

func (s *ServerModel) ConvertProjectToProtos(u []*order.Order) []*pb.Project {
	arr := make([]*pb.Project, 0)
	for _, a := range u {
		arr = append(arr, s.ConvertProjectToProto(a))
	}
	return arr
}

func (s *ServerModel) ConvertOrderToProto(u *order.Order) *pb.Order {
	upb := new(pb.Order)
	if u.ID != uuid.Nil {
		upb.Id = u.ID.String()
	}
	if u.BusinessId != uuid.Nil {
		upb.BusinessId = u.BusinessId.String()
	}
	if u.Status != nil {
		upb.Status = c.ORDER_STATUS(*u.Status)
	}
	if u.ConversationId != uuid.Nil {
		upb.ConversationId = u.ConversationId.String()
	}
	if u.CustomerId != uuid.Nil {
		upb.CustomerId = u.CustomerId.String()
	}
	if u.ServiceId != uuid.Nil {
		upb.ServiceId = u.ServiceId.String()
	}
	if u.CustomerMessage != nil {
		upb.CustomerMessage = *u.CustomerMessage
	}
	if u.CustomerPhone != nil {
		upb.CustomerPhone = *u.CustomerPhone
	}
	if u.CustomerZipcode != nil {
		upb.CustomerZipcode = *u.CustomerZipcode
	}
	if u.CustomerMessage != nil {
		upb.CustomerMessage = *u.CustomerMessage
	}
	if u.ServiceName != nil {
		upb.ServiceName = *u.ServiceName
	}
	if u.CreatedAt != 0 {
		upb.StartDate = u.CreatedAt
	}
	if u.EndDate != nil {
		upb.EndDate = *u.EndDate
	}
	if u.CustomerAvatar != nil {
		upb.Image = *u.CustomerAvatar
	}
	if u.Fee != nil {
		upb.Fee = *u.Fee
	}
	if u.CategoryName != nil {
		upb.ServiceName = *u.CategoryName
	}
	if u.BusinessName != nil {
		upb.BusinessName = *u.BusinessName
	}
	if u.BusinessBanner != nil {
		upb.BusinessBanner = *u.BusinessBanner
	}
	if u.BusinessLogo != nil {
		upb.BusinessLogo = *u.BusinessLogo
	}
	if u.CustomerName != nil {
		upb.CustomerName = *u.CustomerName
	}
	if u.IsReviewed != nil {
		upb.IsReviewed = *u.IsReviewed
	}
	if u.CategoryId != nil {
		upb.CategoryId = *u.CategoryId
	}
	if u.CustomerMail != nil {
		upb.CustomerMail = *u.CustomerMail
	}
	if u.HandymanMail != nil {
		upb.HandymanMail = *u.HandymanMail
	}
	return upb
}

func (s *ServerModel) GetOrderById(ctx context.Context, id interface{}) (*order.Order, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetOrderById))
	defer span.End()

	idd, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	o, err := s.Repo.SelectOrder(ctx, &order.Search{
		Order: order.Order{
			BaseModel: database.BaseModel{ID: idd},
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return o, nil
}

func (s *ServerModel) ConvertOrderToProtos(u []*order.Order) []*pb.Order {
	arr := make([]*pb.Order, 0)
	for _, a := range u {
		arr = append(arr, s.ConvertOrderToProto(a))
	}
	return arr
}

func (s *ServerModel) ListOrders(ctx context.Context, search *order.Search) ([]*order.Order, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListOrders))
	defer span.End()

	err := s.Repo.UpdateOrderStatusIfExpireTime(ctx)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	odrs, err := s.Repo.ListOrders(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return odrs, nil
}

func (s *ServerModel) TotalOrders(ctx context.Context, search *order.Search) (*int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.TotalOrders))
	defer span.End()

	total, err := s.Repo.TotalOrder(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return total, nil
}

func (s *ServerModel) CreateOrder(ctx context.Context, uid, sid interface{}, zipcode *string, phone *string, message *string) (*order.Order, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CreateOrder))
	defer span.End()

	cid, err := lib.ToUUID(uid)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err

	}

	ser, err := s.GetServiceById(ctx, sid)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	if _, err = s.Repo.SelectOrder(ctx, &order.Search{
		Order: order.Order{
			CustomerId:      cid,
			ServiceId:       ser.ID,
			Status:          utils.Int32Ptr(int32(c.ORDER_STATUS_PENDING)),
			CustomerZipcode: zipcode,
		},
	}); err == nil {
		err = xerrors.Errorf("%w", e.ErrAlreadyOrdered)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	oid := uuid.New()
	now := time.Now()
	ord, err := s.Repo.InsertOrder(ctx, &order.Order{
		BaseModel:       database.BaseModel{ID: oid},
		CustomerId:      cid,
		BusinessId:      ser.BusinessId,
		ConversationId:  oid,
		ServiceId:       ser.ID,
		StartDate:       utils.Int64Ptr(now.UnixMilli()),
		EndDate:         utils.Int64Ptr(now.Add(THREE_DAYS).UnixMilli()),
		Status:          utils.Int32Ptr(int32(c.ORDER_STATUS_PENDING)),
		CustomerZipcode: zipcode,
		CustomerMessage: message,
		CustomerPhone:   phone,
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return ord, nil
}

func (s *ServerModel) CreateOrderV2(ctx context.Context, uid, bid, cid interface{}, zipcode *string, phone *string, message *string, customerName *string) (*order.Order, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CreateOrderV2))
	defer span.End()

	uidd, err := lib.ToUUID(uid)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err

	}
	bidd, err := lib.ToUUID(bid)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err

	}
	cidd, err := lib.ToUUID(cid)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err

	}
	ser, err := s.Repo.SelectService(ctx, &service.Search{
		Service: service.Service{
			BusinessId: bidd,
			CategoryId: cidd,
			Status:     utils.Int32Ptr(int32(c.SERVICE_STATUS_SERVICE_ACTIVE)),
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	if _, err = s.Repo.SelectOrder(ctx, &order.Search{
		Order: order.Order{
			CustomerId:      uidd,
			ServiceId:       ser.ID,
			Status:          utils.Int32Ptr(int32(c.ORDER_STATUS_PENDING)),
			CustomerZipcode: zipcode,
		},
	}); err == nil {
		err = xerrors.Errorf("%w", e.ErrAlreadyOrdered)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	now := time.Now()
	ord, err := s.Repo.InsertOrder(ctx, &order.Order{
		CustomerId:      uidd,
		BusinessId:      ser.BusinessId,
		ServiceId:       ser.ID,
		StartDate:       utils.Int64Ptr(now.UnixMilli()),
		EndDate:         utils.Int64Ptr(now.Add(THREE_DAYS).UnixMilli()),
		Status:          utils.Int32Ptr(int32(c.ORDER_STATUS_PENDING)),
		CustomerZipcode: zipcode,
		CustomerMessage: message,
		CustomerPhone:   phone,
		CustomerName:    customerName,
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return ord, nil
}

func (s *ServerModel) GetCurrentOrderCount(ctx context.Context, uid interface{}, zipcode *string) (*int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetCurrentOrderCount))
	defer span.End()

	cid, err := lib.ToUUID(uid)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	count, err := s.Repo.TotalOrder(ctx, &order.Search{
		Order: order.Order{
			CustomerId:      cid,
			CustomerZipcode: zipcode,
			Status:          utils.Int32Ptr(int32(c.ORDER_STATUS_PENDING)),
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return count, nil

}

func (s *ServerModel) UpdateOrderStatusById(ctx context.Context, orderId interface{}, status c.ORDER_STATUS) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UpdateOrderStatusById))
	defer span.End()

	uid, err := lib.ToUUID(orderId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	err = s.Repo.UpdateOrder(ctx, &order.Search{
		Order: order.Order{
			BaseModel: database.BaseModel{ID: uid},
		},
	}, &order.Order{
		Status: utils.Int32Ptr(int32(status)),
	})

	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}

	return nil
}

func (s *ServerModel) UpdateOrdersStatusByUser(ctx context.Context, userId interface{}, status c.ORDER_STATUS) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UpdateOrderStatusById))
	defer span.End()

	uid, err := lib.ToUUID(userId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	err = s.Repo.UpdateOrdersCancelledStatusByUser(ctx, &order.Search{
		UserId: uid,
	}, &order.Order{
		Status: utils.Int32Ptr(int32(status)),
	})

	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}

	return nil
}

func (s *ServerModel) CheckPermissionUpdateOrder(idRequest string, idFromOrder uuid.UUID) bool {
	uid := idFromOrder.String()

	return uid == idRequest
}

func (s *ServerModel) ListProjects(ctx context.Context, search *order.Search) ([]*order.Order, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListProjects))
	defer span.End()

	search.DefaultSearchModel.Fields = []string{
		`"orders"."customer_zipcode"`,
		`"categories"."name" as "service_name"`,
		`count(*) as "number_orders"`,
		`"categories"."id" as "service_id"`,
		// ADD CATEGORY IMAGE
	}

	prj, err := s.Repo.ListProjects(ctx, search)

	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return prj, nil
}

func (s *ServerModel) TotalProjects(ctx context.Context, search *order.Search) (*int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.TotalProjects))
	defer span.End()

	total, err := s.Repo.TotalProjects(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return total, nil
}

func (s *ServerModel) CancelProject(ctx context.Context, search *order.Search) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CancelProject))
	defer span.End()

	err := s.Repo.CancelProject(ctx, search, &order.Order{
		Status: utils.Int32Ptr(int32(c.ORDER_STATUS_CANCELED)),
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}

func (s *ServerModel) ListBusinessesAlreadyOrdered(ctx context.Context, search *order.Search) ([]*order.Order, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListBusinessesAlreadyOrdered))
	defer span.End()

	ords, err := s.Repo.ListBusinessesAlreadyOrdered(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return ords, nil
}
