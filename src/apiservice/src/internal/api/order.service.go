package api

import (
	"context"
	"strconv"

	"github.com/aqaurius6666/apiservice/src/internal/db/business"
	"github.com/aqaurius6666/apiservice/src/internal/db/group"
	"github.com/aqaurius6666/apiservice/src/internal/db/order"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/lib/unleash"
	"github.com/aqaurius6666/apiservice/src/internal/lib/validate"
	"github.com/aqaurius6666/apiservice/src/internal/model"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/apiservice/src/internal/var/e"
	"github.com/aqaurius6666/apiservice/src/pb"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/aqaurius6666/go-utils/utils"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

type OrderService struct {
	Model  model.Server
	Logger *logrus.Logger
}

func (s OrderService) CreateOrder(ctx context.Context, req *pb.OrdersPostRequest) (*pb.OrdersPostResponse_Data, error) {

	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CreateOrder))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId", "CategoryId", "Zipcode", "BusinessIds"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	currOrder, err := s.Model.GetCurrentOrderCount(ctx, req.XUserId, &req.Zipcode)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	maxOrder, err := strconv.Atoi(unleash.GetVariant("apiservice.config.max-order"))
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	usr, err := s.Model.GetUserById(ctx, req.XUserId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	customerName := utils.StrVal(usr.FirstName) + " " + utils.StrVal(usr.LastName)

	if int(*currOrder)+len(req.BusinessIds) > maxOrder {
		err = xerrors.Errorf("%w", e.ErrExceedMaxOrders)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	cate, err := s.Model.GetCategoryById(ctx, req.CategoryId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	for _, bid := range req.BusinessIds {
		_, err := s.Model.CreateOrderV2(ctx, req.XUserId, bid, req.CategoryId, &req.Zipcode, usr.Phone, nil, &customerName)
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err, ctx)
			return nil, err
		}
		go func(handymanId, customerName, categoryName, zipcode string) {
			err := s.Model.SendRequestedNotification(context.TODO(), handymanId, customerName, categoryName, zipcode)
			if err != nil {
				s.Logger.Error(err)
			}
		}(bid, customerName, *cate.Name, req.Zipcode)

	}
	return &pb.OrdersPostResponse_Data{}, nil
}

func (s OrderService) ListOrders(ctx context.Context, req *pb.OrdersGetRequest) (*pb.OrdersGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListOrders))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	limit := lib.ParseInt32Val(req.Limit)
	offset := lib.ParseInt32Val(req.Offset)
	total, err := s.Model.TotalOrders(ctx, &order.Search{
		Order: order.Order{
			Status:          utils.Int32Ptr(int32(req.Status)),
			CustomerZipcode: utils.SafeStrPtr(req.Zipcode),
		},
		UserId:     lib.ParseUUID(req.XUserId),
		CategoryId: lib.ParseUUID(req.ServiceId),
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	orders, err := s.Model.ListOrders(ctx, &order.Search{
		DefaultSearchModel: database.DefaultSearchModel{
			Skip:  int(offset),
			Limit: int(limit),
			Fields: []string{`"orders".*`,
				`"users"."avatar_url" as customer_avatar`,
				`"groups"."fee"`,
				`"categories"."name" as category_name`,
				`"businesses"."name" as business_name`,
				`"businesses"."logo_url" as business_logo`,
				`"businesses"."banner_url" as business_banner`,
				`"businesses"."mail" as handyman_mail`,
				`"users"."mail" as customer_mail`,
				`"services"."category_id"`,
			},
		},
		Order: order.Order{
			Status:          utils.Int32Ptr(int32(req.Status)),
			CustomerZipcode: utils.SafeStrPtr(req.Zipcode),
		},
		UserId:     lib.ParseUUID(req.XUserId),
		CategoryId: lib.ParseUUID(req.ServiceId),
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &pb.OrdersGetResponse_Data{
		Pagination: lib.Pagination(offset, limit, total),
		Result:     s.Model.ConvertOrderToProtos(orders),
	}, nil
}

func (s OrderService) UpdateOrderConnectedStatus(ctx context.Context, req *pb.UpdateOrderStatusPostRequest) (*pb.UpdateOrderStatusPostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UpdateOrderConnectedStatus))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId", "OrderId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	ord, err := s.Model.GetOrderById(ctx, req.OrderId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	if !s.Model.CheckPermissionUpdateOrder(req.XUserId, ord.BusinessId) {
		err = xerrors.Errorf("%w", e.ErrNoPermission)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	if *ord.Status != int32(c.ORDER_STATUS_PENDING) {
		err = xerrors.Errorf("%w", e.ErrInvalidOrderStatus)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	convId, err := s.Model.NewConversation(ctx, ord.BusinessId, ord.CustomerId, ord.ID)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	sv, err := s.Model.GetServiceById(ctx, ord.ServiceId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	gr, err := s.Model.SelectGroup(ctx, &group.Search{CategoryId: sv.CategoryId})
	if err != nil {
		// Category not in group
		if err = s.Model.UpdateOrderById(ctx, ord.ID, &order.Order{
			ConversationId: convId,
			Status:         utils.Int32Ptr(int32(c.ORDER_STATUS_CONNECTED)),
		}); err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err)
			return nil, err
		}

		return &pb.UpdateOrderStatusPostResponse_Data{}, nil
	}

	bus, err := s.Model.GetBusinessById(ctx, req.XUserId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	if utils.Int32Val(bus.FreeContact) > 0 {
		err = s.Model.InsertTransaction(ctx, ord.ID, req.XUserId, *gr.Fee, true)
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err, ctx)
			return nil, err
		}

		freeContact := utils.Int32Val(bus.FreeContact) - 1
		err = s.Model.UpdateBusiness(ctx, req.XUserId, &business.Business{
			FreeContact: &freeContact,
		})
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err, ctx)
			return nil, err
		}
	} else {
		err = s.Model.InsertTransaction(ctx, ord.ID, req.XUserId, *gr.Fee, false)
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err, ctx)
			return nil, err
		}
	}

	err = s.Model.UpdateOrderById(ctx, ord.ID, &order.Order{
		ConversationId: convId,
		Status:         utils.Int32Ptr(int32(c.ORDER_STATUS_CONNECTED)),
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	go func(customerId, businessName string) {
		err := s.Model.SendConnectNotification(context.TODO(), customerId, businessName, convId.String())
		if err != nil {
			s.Logger.Error(err)
		}
	}(ord.CustomerId.String(), *bus.Name)

	return &pb.UpdateOrderStatusPostResponse_Data{}, nil
}

func (s OrderService) UpdateAllOrderConnectedStatus(ctx context.Context, req *pb.UpdateAllOrderStatusPostRequest) (*pb.UpdateAllOrderStatusPostResponse, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UpdateAllOrderConnectedStatus))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	bid, err := lib.ToUUID(req.XUserId)
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}

	orders, err := s.Model.ListOrders(ctx, &order.Search{
		Order: order.Order{
			BusinessId: bid,
			Status:     utils.Int32Ptr(int32(c.ORDER_STATUS_PENDING)),
		},
		DefaultSearchModel: database.DefaultSearchModel{
			Fields: []string{`"orders".*`},
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	for _, ord := range orders {
		if !s.Model.CheckPermissionUpdateOrder(req.XUserId, ord.BusinessId) {
			err = xerrors.Errorf("%w", e.ErrNoPermission)
			lib.RecordError(span, err, ctx)
			return nil, err
		}

		convId, err := s.Model.NewConversation(ctx, ord.BusinessId, ord.CustomerId, ord.ID)
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err, ctx)
			return nil, err
		}

		sv, err := s.Model.GetServiceById(ctx, ord.ServiceId)
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err, ctx)
			return nil, err
		}

		gr, err := s.Model.SelectGroup(ctx, &group.Search{CategoryId: sv.CategoryId})
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err, ctx)
			return nil, err
		}

		bus, err := s.Model.GetBusinessById(ctx, req.XUserId)
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err, ctx)
			return nil, err
		}

		if utils.Int32Val(bus.FreeContact) > 0 {
			err = s.Model.InsertTransaction(ctx, ord.ID, req.XUserId, *gr.Fee, true)
			if err != nil {
				err = xerrors.Errorf("%w", err)
				lib.RecordError(span, err, ctx)
				return nil, err
			}

			freeContact := utils.Int32Val(bus.FreeContact) - 1
			err = s.Model.UpdateBusiness(ctx, req.XUserId, &business.Business{
				FreeContact: &freeContact,
			})
			if err != nil {
				err = xerrors.Errorf("%w", err)
				lib.RecordError(span, err, ctx)
				return nil, err
			}
		} else {
			err = s.Model.InsertTransaction(ctx, ord.ID, req.XUserId, *gr.Fee, false)
			if err != nil {
				err = xerrors.Errorf("%w", err)
				lib.RecordError(span, err, ctx)
				return nil, err
			}
		}

		err = s.Model.UpdateOrderById(ctx, ord.ID, &order.Order{
			ConversationId: convId,
			Status:         utils.Int32Ptr(int32(c.ORDER_STATUS_CONNECTED)),
		})
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err, ctx)
			return nil, err
		}
		go func(customerId, businessName string) {
			err := s.Model.SendConnectNotification(context.TODO(), customerId, businessName, convId.String())
			if err != nil {
				s.Logger.Error(err)
			}
		}(ord.CustomerId.String(), *bus.Name)
	}

	return &pb.UpdateAllOrderStatusPostResponse{}, nil
}

func (s OrderService) UpdateOrderRejectedStatus(ctx context.Context, req *pb.UpdateOrderStatusPostRequest) (*pb.UpdateOrderStatusPostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UpdateOrderConnectedStatus))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId", "OrderId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	bus, err := s.Model.GetBusinessById(ctx, req.XUserId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	order, err := s.Model.GetOrderById(ctx, req.OrderId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	if !s.Model.CheckPermissionUpdateOrder(req.XUserId, order.BusinessId) {
		err = xerrors.Errorf("%w", e.ErrNoPermission)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	if !(*order.Status == int32(c.ORDER_STATUS_PENDING) || *order.Status == int32(c.ORDER_STATUS_CONNECTED)) {
		err = xerrors.Errorf("%w", e.ErrInvalidOrderStatus)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	err = s.Model.UpdateOrderStatusById(ctx, req.OrderId, c.ORDER_STATUS_REJECTED)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	err = s.Model.CloseConversation(ctx, order.ConversationId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	go func(customerId, businessName, businessId string) {
		err := s.Model.SendRejectNotification(context.TODO(), customerId, businessName, businessId)
		if err != nil {
			s.Logger.Error(err)
		}
	}(order.CustomerId.String(), *bus.Name, bus.ID.String())

	return &pb.UpdateOrderStatusPostResponse_Data{}, nil
}

func (s OrderService) UpdateOrderCompletedStatus(ctx context.Context, req *pb.UpdateOrderStatusPostRequest) (*pb.UpdateOrderStatusPostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UpdateOrderConnectedStatus))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId", "OrderId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	order, err := s.Model.GetOrderById(ctx, req.OrderId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	if req.XRole == c.ROLE_HANDYMAN && !s.Model.CheckPermissionUpdateOrder(req.XUserId, order.BusinessId) {
		err = xerrors.Errorf("%w", e.ErrNoPermission)
		lib.RecordError(span, err, ctx)
		return nil, err
	} else if req.XRole == c.ROLE_CUSTOMER && !s.Model.CheckPermissionUpdateOrder(req.XUserId, order.CustomerId) {
		err = xerrors.Errorf("%w", e.ErrNoPermission)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	bus, err := s.Model.GetBusinessById(ctx, order.BusinessId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	if *order.Status != int32(c.ORDER_STATUS_CONNECTED) {
		err = xerrors.Errorf("%w", e.ErrInvalidOrderStatus)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	err = s.Model.UpdateOrderStatusById(ctx, req.OrderId, c.ORDER_STATUS_COMPLETED)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	err = s.Model.CloseConversation(ctx, order.ConversationId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	go func(customerId, businessName, businessId string) {
		err := s.Model.SendCompleteNotification(context.TODO(), customerId, businessName, businessId)
		if err != nil {
			s.Logger.Error(err)
		}
	}(order.CustomerId.String(), *bus.Name, bus.ID.String())

	return &pb.UpdateOrderStatusPostResponse_Data{}, nil
}

func (s OrderService) UpdateOrderCancelledStatus(ctx context.Context, req *pb.UpdateOrderStatusPostRequest) (*pb.UpdateOrderStatusPostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UpdateOrderConnectedStatus))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId", "OrderId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	order, err := s.Model.GetOrderById(ctx, req.OrderId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	if !s.Model.CheckPermissionUpdateOrder(req.XUserId, order.CustomerId) {
		err = xerrors.Errorf("%w", e.ErrNoPermission)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	if *order.Status != int32(c.ORDER_STATUS_PENDING) {
		err = xerrors.Errorf("%w", e.ErrInvalidOrderStatus)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	err = s.Model.UpdateOrderStatusById(ctx, req.OrderId, c.ORDER_STATUS_CANCELED)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	service, err := s.Model.GetServiceById(ctx, order.ServiceId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	go func(handymanId, customerName, categoryName, zipcode string) {
		err := s.Model.SendCancelNotification(context.TODO(), handymanId, customerName, categoryName, zipcode)
		if err != nil {
			s.Logger.Error(err)
		}
	}(order.BusinessId.String(), *order.CustomerName, *service.Category.Name, *order.CustomerZipcode)

	return &pb.UpdateOrderStatusPostResponse_Data{}, nil
}

func (s OrderService) GetProjects(ctx context.Context, req *pb.UserProjectsGetRequest) (*pb.UserProjectsGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetProjects))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	limit := lib.ParseInt32Val(req.Limit)
	offset := lib.ParseInt32Val(req.Offset)

	uid, err := lib.ToUUID(req.XUserId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	prj, err := s.Model.ListProjects(ctx, &order.Search{
		DefaultSearchModel: database.DefaultSearchModel{
			Limit: int(limit),
			Skip:  int(offset),
		},
		Order: order.Order{
			CustomerId: uid,
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	total, err := s.Model.TotalProjects(ctx, &order.Search{
		Order: order.Order{
			CustomerId: uid,
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.UserProjectsGetResponse_Data{
		Result:     s.Model.ConvertProjectToProtos(prj),
		Pagination: lib.Pagination(offset, limit, total),
	}, nil
}

func (s OrderService) CancelProject(ctx context.Context, req *pb.CancelProjectPostRequest) (*pb.CancelProjectPostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CancelProject))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId", "Zipcode", "CategoryId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	uid, err := lib.ToUUID(req.XUserId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	sid, err := lib.ToUUID(req.CategoryId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	err = s.Model.CancelProject(ctx, &order.Search{
		Order: order.Order{
			CustomerZipcode: &req.Zipcode,
			ServiceId:       sid,
			CustomerId:      uid,
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.CancelProjectPostResponse_Data{}, nil
}

func (s OrderService) BusinessesAlreadyOrdered(ctx context.Context, req *pb.BusinessesAlreadyOrderedGetRequest) (*pb.BusinessesAlreadyOrderedGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CancelProject))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId", "Zipcode", "CategoryId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	uid, err := lib.ToUUID(req.XUserId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	cid, err := lib.ToUUID(req.CategoryId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	orders, err := s.Model.ListBusinessesAlreadyOrdered(ctx, &order.Search{
		CategoryId: cid,
		Order: order.Order{
			CustomerId:      uid,
			CustomerZipcode: &req.Zipcode,
		},
		DefaultSearchModel: database.DefaultSearchModel{Fields: []string{`"orders"."business_id"`}},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	result := make([]string, 0)
	for _, s := range orders {
		result = append(result, s.BusinessId.String())
	}

	return &pb.BusinessesAlreadyOrderedGetResponse_Data{
		BusinessId: result,
	}, nil
}
