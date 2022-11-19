package api

import (
	"context"
	"io"

	"github.com/aqaurius6666/apiservice/src/internal/db/advertise_order"
	"github.com/aqaurius6666/apiservice/src/internal/db/advertise_package"
	"github.com/aqaurius6666/apiservice/src/internal/db/advertise_transaction"
	"github.com/aqaurius6666/apiservice/src/internal/db/business"
	"github.com/aqaurius6666/apiservice/src/internal/db/feedback"
	"github.com/aqaurius6666/apiservice/src/internal/db/payment"
	"github.com/aqaurius6666/apiservice/src/internal/db/service"
	"github.com/aqaurius6666/apiservice/src/internal/db/transaction"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/lib/validate"
	"github.com/aqaurius6666/apiservice/src/internal/model"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/apiservice/src/internal/var/e"
	"github.com/aqaurius6666/apiservice/src/pb"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/aqaurius6666/go-utils/utils"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

type BusinessService struct {
	Model model.Server
}

func (s BusinessService) GetInterest(ctx context.Context, req *pb.BusinessInterestGetRequest) (*pb.BusinessInterestGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetInterest))
	defer span.End()
	buss, err := s.Model.ListBusinessesWithRating(ctx, &business.Search{
		DefaultSearchModel: database.DefaultSearchModel{
			OrderBy:   `"rating"."rate"`,
			OrderType: "DESC",
			Limit:     20,
			Fields: []string{`"businesses"."id"`, "name", "logo_url", "banner_url", "contact_id",
				"website", "description", "services", `"rating"."rate"`, `"rating"."review"`},
		},
		Business: business.Business{
			Status: utils.Int32Ptr(int32(c.ACCOUNT_STATUS_ACTIVE)),
		},
	})
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}
	ret := make([]*pb.BusinessRating, 0)
	for _, b := range buss {
		ret = append(ret, &pb.BusinessRating{
			Business: s.Model.ConvertBusinessToProto(b),
			Rating: s.Model.ConvertRatingToProto(&feedback.Feedback{
				Rate:    b.Rate,
				Review:  b.Review,
				Request: b.Request,
			}),
		})
	}
	return &pb.BusinessInterestGetResponse_Data{
		Result: ret,
	}, nil
}

func (s BusinessService) List(ctx context.Context, req *pb.BusinessesGetRequest) (*pb.BusinessesGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.List))
	defer span.End()

	limit := lib.ParseInt32Val(req.Limit)
	offset := lib.ParseInt32Val(req.Offset)
	categoryId := lib.ParseUUID(req.CategoryId)
	search := &business.Search{
		DefaultSearchModel: database.DefaultSearchModel{
			Skip:  int(offset),
			Limit: int(limit),
			Fields: []string{
				`"businesses"."id"`,
				`"businesses"."name"`,
				`"businesses"."logo_url"`,
				`"businesses"."banner_url"`,
				`"businesses"."contact_id"`,
				`"businesses"."website"`,
				`"businesses"."description"`,
				`"businesses"."services"`,
				"rate",
				"review",
				"request",
				"start_date",
			},
			OrderBy:   `"rating"."rate" DESC,"rating"."review"`,
			OrderType: "DESC",
		},
		Business: business.Business{
			Phone:  utils.SafeStrPtr(req.Phone),
			Status: utils.Int32Ptr(int32(c.ACCOUNT_STATUS_ACTIVE)),
		},
		Zipcode:    utils.SafeStrPtr(req.Zipcode),
		CategoryId: categoryId,
		Query:      req.Query,
	}

	buss, err := s.Model.ListBusinessesWithRating(ctx, search)
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}
	total, err := s.Model.TotalBusinesses(ctx, search)
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}

	ret := make([]*pb.BusinessRating, 0)
	for _, b := range buss {
		ret = append(ret, &pb.BusinessRating{
			Business: s.Model.ConvertBusinessToProto(b),
			Rating: s.Model.ConvertRatingToProto(&feedback.Feedback{
				Rate:    b.Rate,
				Review:  b.Review,
				Request: b.Request,
			}),
		})
	}
	return &pb.BusinessesGetResponse_Data{
		Result:     ret,
		Pagination: lib.Pagination(offset, limit, total),
	}, nil
}

func (s BusinessService) GetNear(ctx context.Context, req *pb.BusinessNearGetRequest) (*pb.BusinessNearGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetNear))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	usr, err := s.Model.GetUserById(ctx, req.XUserId)
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}
	buss, err := s.Model.GetNearBussiness(ctx, usr.Contact.Zipcode)
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}

	ret := make([]*pb.BusinessRating, 0)
	for _, b := range buss {
		ret = append(ret, &pb.BusinessRating{
			Business: s.Model.ConvertBusinessToProto(b),
			Rating: s.Model.ConvertRatingToProto(&feedback.Feedback{
				Rate:   b.Rate,
				Review: b.Review,
			}),
		})
	}
	return &pb.BusinessNearGetResponse_Data{
		Result: ret,
	}, nil
}
func (s BusinessService) GetServices(ctx context.Context, req *pb.BusinessServiceGetRequest) (*pb.BusinessServiceGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetServices))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "Id"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	ss, err := s.Model.GetServicesByBusiness(ctx, req.Id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &pb.BusinessServiceGetResponse_Data{
		Result: s.Model.ConvertServiceToProtos(ss),
	}, nil
}

func (s BusinessService) UpdateServices(ctx context.Context, req *pb.BusinessServicesPutRequest) (*pb.BusinessServicesPutResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UpdateServices))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "Id", "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	if req.Id != req.XUserId {
		err := xerrors.Errorf("%w", e.ErrNoPermission)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	uids := make([]uuid.UUID, 0)
	for _, c := range req.CategoryIds {
		uid, err := uuid.Parse(c)
		if err != nil {
			err = xerrors.Errorf("%w", e.ErrIdInvalidFormat)
			lib.RecordError(span, err, ctx)
			return nil, err
		}
		uids = append(uids, uid)
	}
	ss, err := s.Model.UpdateServices(ctx, req.Id, uids)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.BusinessServicesPutResponse_Data{
		Result: s.Model.ConvertServiceToProtos(ss),
	}, nil
}

func (s BusinessService) CreateBusiness(ctx context.Context, req *pb.BusinessPostRequest) (*pb.BusinessPostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CreateBusiness))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "Mail", "PublicKey", "EncryptedPrivateKey"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	if err := req.Validate(); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	var refCode = utils.SafeStrPtr(req.RefCode)
	_, err := s.Model.GetBusinessByReferralCode(ctx, refCode)
	if err != nil {
		err := xerrors.Errorf("%w", e.RefCodeNotFound)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	otpId, err := s.Model.Register(ctx, req.Mail, req.Phone, req.PublicKey, req.EncryptedPrivateKey, c.ROLE_HANDYMAN, req.RefCode)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.BusinessPostResponse_Data{
		OtpId: otpId,
		Mail:  req.Mail,
		Phone: req.Phone,
	}, nil
}

func (s BusinessService) GetById(ctx context.Context, req *pb.BusinessGetRequest) (*pb.BusinessGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetById))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "Id"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	b, err := s.Model.GetBusinessById(ctx, req.Id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.BusinessGetResponse_Data{
		Business: s.Model.ConvertBusinessToProto(b),
	}, nil
}

func (s BusinessService) Update(ctx context.Context, req *pb.BusinessPutRequest) (*pb.BusinessPutResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.Update))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "Id", "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	if req.Id != req.XUserId {
		err := xerrors.Errorf("%w", e.ErrNoPermission)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	if err := req.Validate(); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	b := &business.Business{
		Name:        utils.SafeStrPtr(req.Name),
		Phone:       utils.SafeStrPtr(req.Phone),
		LogoUrl:     utils.SafeStrPtr(req.LogoUrl),
		BannerUrl:   utils.SafeStrPtr(req.BannerUrl),
		Website:     utils.SafeStrPtr(req.Website),
		Description: utils.SafeStrPtr(req.Description),
		Zipcodes:    req.Zipcodes,
	}
	if err := validate.Validate(b); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	if err := s.Model.UpdateBusiness(ctx, req.Id, b); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.BusinessPutResponse_Data{
		Business: s.Model.ConvertBusinessToProto(b),
	}, nil
}
func (s BusinessService) SetupPaymentMethod(ctx context.Context, req *pb.BusinessPaymentMethodSetupPostRequest) (*pb.BusinessPaymentMethodSetupPostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.SetupPaymentMethod))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	bus, err := s.Model.GetPaymentByBusinessId(ctx, req.XUserId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	if bus.CustomerId == nil {
		err = xerrors.Errorf("%w", e.ErrBodyInvalid)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	clientSecret, err := s.Model.SetupPaymentMethod(ctx, bus.CustomerId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &pb.BusinessPaymentMethodSetupPostResponse_Data{
		ClientSecret: *clientSecret,
	}, nil
}

func (s BusinessService) DeletePaymentMethod(ctx context.Context, req *pb.BusinessPaymentMethodDeletePostRequest) (*pb.BusinessPaymentMethodDeletePostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.DeletePaymentMethod))
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

	bus, err := s.Model.GetPaymentByBusinessId(ctx, req.XUserId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	if bus.CustomerId == nil || bus.PaymentMethodId == nil {
		err = xerrors.Errorf("%w", e.ErrBodyInvalid)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	var emptyStr string = ""

	err = s.Model.UpdatePayment(ctx, &payment.Search{
		Payment: payment.Payment{
			BusinessId: bid,
		},
	}, &payment.Payment{
		PaymentMethodId: &emptyStr,
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.BusinessPaymentMethodDeletePostResponse_Data{}, nil
}

func (s BusinessService) GetRatings(ctx context.Context, req *pb.BusinessRatingGetRequest) (*pb.BusinessRatingGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetRatings))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "Id"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	fb, err := s.Model.ListRatingByBusiness(ctx, req.Id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &pb.BusinessRatingGetResponse_Data{
		Rate: s.Model.ConvertRatingToProtos(fb),
	}, nil
}

func (s BusinessService) GetFeedbacks(ctx context.Context, req *pb.BusinessFeedbacksGetRequest) (*pb.BusinessFeedbacksGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetFeedbacks))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "Id"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	uid, err := lib.ToUUID(req.Id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	limit := lib.ParseInt32Val(req.Limit)
	offset := lib.ParseInt32Val(req.Offset)

	search := &feedback.Search{
		DefaultSearchModel: database.DefaultSearchModel{
			Skip:  int(offset),
			Limit: int(limit),
			Fields: []string{
				`CONCAT("first_name", ' ', "last_name") AS "customer_name"`,
				"rate",
				"comment",
				`"feedbacks"."created_at"`,
				"avatar_url",
				`"categories"."name" AS "service_name"`,
			},
		},
		Feedback: feedback.Feedback{
			BusinessId: uid,
		},
	}

	search1 := &feedback.Search{
		Feedback: feedback.Feedback{
			BusinessId: uid,
		},
	}

	fb, err := s.Model.ListFeedbackByBusiness(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	total, err := s.Model.TotalFeedbackByBusiness(ctx, search1)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.BusinessFeedbacksGetResponse_Data{
		Result:     s.Model.ConvertFeedbackToProtos(fb),
		Pagination: lib.Pagination(offset, limit, total),
	}, nil
}

func (s BusinessService) GetPaymentMethod(ctx context.Context, req *pb.BusinessPaymentMethodGetRequest) (*pb.BusinessPaymentMethodGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetPaymentMethod))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	// get and verify business info
	b, err := s.Model.GetBusinessById(ctx, req.XUserId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	// check if have payment or not
	pay, err := s.Model.GetPaymentByBusinessId(ctx, req.XUserId)
	if err == nil {
		return &pb.BusinessPaymentMethodGetResponse_Data{
			Payment: s.Model.ConvertPaymentToProto(pay),
		}, nil
	}

	// create key and put into payment table
	payment, err := s.Model.NewCustomer(ctx, req.XUserId, *b.Mail)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.BusinessPaymentMethodGetResponse_Data{
		Payment: s.Model.ConvertPaymentToProto(payment),
	}, nil
}

func (s BusinessService) PostPaymentMethod(ctx context.Context, req *pb.BusinessPaymentMethodPostRequest) (*pb.BusinessPaymentMethodPostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.PostPaymentMethod))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId", "PaymentMethodId"); !ok {
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

	// get and verify business info
	_, err = s.Model.GetBusinessById(ctx, req.XUserId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	// check if have payment or not. If not, return err
	_, err = s.Model.GetPaymentByBusinessId(ctx, req.XUserId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	// update payment_method_id into table payment

	err = s.Model.UpdatePayment(ctx, &payment.Search{
		Payment: payment.Payment{
			BusinessId: uid,
		},
	}, &payment.Payment{
		PaymentMethodId: &req.PaymentMethodId,
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &pb.BusinessPaymentMethodPostResponse_Data{}, nil
}

func (s BusinessService) ListTransactions(ctx context.Context, req *pb.BusinessTransactionsGetRequest) (*pb.BusinessTransactionsGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListTransactions))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId"); !ok {
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

	limit := lib.ParseInt32Val(req.Limit)
	offset := lib.ParseInt32Val(req.Offset)

	total, err := s.Model.TotalTransaction(ctx, &transaction.Search{
		Query: utils.Int32Ptr(int32(req.Query)),
		Transaction: transaction.Transaction{
			BusinessId: uid,
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	totalFee, err := s.Model.TotalFee(ctx, &transaction.Search{
		Query: utils.Int32Ptr(int32(req.Query)),
		Transaction: transaction.Transaction{
			BusinessId: uid,
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	transactions, err := s.Model.GetTransactions(ctx, &transaction.Search{
		DefaultSearchModel: database.DefaultSearchModel{
			Skip:  int(offset),
			Limit: int(limit),
			Fields: []string{
				`"orders"."start_date"`,
				`"orders"."end_date"`,
				`"categories"."name" as "service_name"`,
				`"orders"."customer_zipcode"`,
				`"transactions"."fee"`,
				`"orders"."status"`,
				`"users"."avatar_url" as customer_avatar`,
			},
		},
		Query: utils.Int32Ptr(int32(req.Query)),
		Transaction: transaction.Transaction{
			BusinessId: uid,
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.BusinessTransactionsGetResponse_Data{
		Pagination: lib.Pagination(offset, limit, total),
		Result:     s.Model.ConvertTransactionToProtos(transactions),
		TotalFee:   *lib.CentToUsd(totalFee),
	}, nil
}

func (s BusinessService) TransactionsExport(ctx context.Context, req *pb.BusinessTransactionsGetRequest) (io.Reader, int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.TransactionsExport))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, 0, err
	}

	uid, err := lib.ToUUID(req.XUserId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, 0, err
	}

	totalFee, err := s.Model.TotalFee(ctx, &transaction.Search{
		Query: utils.Int32Ptr(int32(req.Query)),
		Transaction: transaction.Transaction{
			BusinessId: uid,
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, 0, err
	}

	transactions, err := s.Model.GetTransactions(ctx, &transaction.Search{
		DefaultSearchModel: database.DefaultSearchModel{
			Fields: []string{
				`"orders"."start_date"`,
				`"orders"."end_date"`,
				`"categories"."name" as "service_name"`,
				`"orders"."customer_zipcode"`,
				`"transactions"."fee"`,
				`"orders"."status"`,
				`"users"."avatar_url" as customer_avatar`,
			},
		},
		Query: utils.Int32Ptr(int32(req.Query)),
		Transaction: transaction.Transaction{
			BusinessId: uid,
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, 0, err
	}

	buf, len, err := s.Model.ExportTransactions(ctx, transactions, lib.CentToUsd(totalFee))
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, 0, err
	}

	return buf, len, nil
}

func (s *BusinessService) GetAdvertisePackages(ctx context.Context, req *pb.AdvertiseGetRequest) (*pb.AdvertiseGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetAdvertisePackages))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	limit := lib.ParseInt32Val(req.Limit)
	offset := lib.ParseInt32Val(req.Offset)

	total, err := s.Model.TotalPackage(ctx, &advertise_package.Search{
		DefaultSearchModel: database.DefaultSearchModel{
			Limit: int(limit),
			Skip:  int(offset),
			Fields: []string{`"advertise_packages"."id"`,
				`"advertise_packages"."name"`,
				`"advertise_packages"."price"`,
				`"advertise_packages"."banner_url"`,
				`"advertise_packages"."description"`,
				`array_remove(array_agg("categories"."name"), null) as "service_name"`,
				`array_remove(array_agg("categories"."id"), null) as "service_id"`},
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	res, err := s.Model.ListPackages(ctx, &advertise_package.Search{
		DefaultSearchModel: database.DefaultSearchModel{
			Limit: int(limit),
			Skip:  int(offset),
			Fields: []string{`"advertise_packages"."id"`,
				`"advertise_packages"."name"`,
				`"advertise_packages"."price"`,
				`"advertise_packages"."banner_url"`,
				`"advertise_packages"."description"`,
				`array_remove(array_agg("categories"."name"), null) as "service_name"`,
				`array_remove(array_agg("categories"."id"), null) as "service_id"`},
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.AdvertiseGetResponse_Data{
		Result:     s.Model.ConvertPackageToProtos(res),
		Pagination: lib.Pagination(offset, limit, total),
	}, nil
}

func (s *BusinessService) GetAdvertiseDetail(ctx context.Context, req *pb.AdvertiseDetailGetRequest) (*pb.AdvertiseDetailGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetAdvertiseDetail))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	businessId, err := lib.ToUUID(req.XUserId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	advertiseId, err := lib.ToUUID(req.Id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	limit := lib.ParseInt32Val(req.Limit)
	offset := lib.ParseInt32Val(req.Offset)

	total, err := s.Model.TotalAdvertiseDetail(ctx, &advertise_package.Search{
		DefaultSearchModel: database.DefaultSearchModel{
			Limit: int(limit),
			Skip:  int(offset),
			Fields: []string{`"advertise_packages"."id"`,
				`"advertise_packages"."name"`,
				`"advertise_packages"."price"`,
				`"businesses"."zipcodes"`,
				`"advertise_packages"."banner_url"`,
				`"advertise_packages"."description"`,
				`array_remove(array_agg("categories"."name"), null) as "service_name"`,
				`array_remove(array_agg("categories"."id"), null) as "service_id"`},
		},
		AdvertisePackage: advertise_package.AdvertisePackage{
			BaseModel: database.BaseModel{
				ID: advertiseId,
			},
		},
		BusinessId: businessId,
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	res, err := s.Model.ListAdvertiseDetails(ctx, &advertise_package.Search{
		DefaultSearchModel: database.DefaultSearchModel{
			Limit: int(limit),
			Skip:  int(offset),
			Fields: []string{`"advertise_packages"."id"`,
				`"advertise_packages"."name"`,
				`"advertise_packages"."price"`,
				`"businesses"."zipcodes"`,
				`"advertise_packages"."banner_url"`,
				`"advertise_packages"."description"`,
				`array_remove(array_agg("categories"."name"), null) as "service_name"`,
				`array_remove(array_agg("categories"."id"), null) as "service_id"`},
		},
		AdvertisePackage: advertise_package.AdvertisePackage{
			BaseModel: database.BaseModel{
				ID: advertiseId,
			},
		},
		BusinessId: businessId,
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.AdvertiseDetailGetResponse_Data{
		Result:     s.Model.ConvertAdvertiseDetailToProtos(res),
		Pagination: lib.Pagination(offset, limit, total),
	}, nil
}

func (s *BusinessService) GetAdvertiseOrder(ctx context.Context, req *pb.BusinessAdvertiseOrderGetRequest) (*pb.BusinessAdvertiseOrderGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetAdvertiseOrder))
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

	total, err := s.Model.TotalAdvertiseOrder(ctx, &advertise_order.Search{
		AdvertiseOrder: advertise_order.AdvertiseOrder{
			BusinessId: uid,
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	res, err := s.Model.ListAdvertiseOrders(ctx, &advertise_order.Search{
		DefaultSearchModel: database.DefaultSearchModel{
			Limit: int(limit),
			Skip:  int(offset),
			Fields: []string{`"advertise_orders"."id"`,
				`"advertise_orders"."start_date"`,
				`"advertise_orders"."end_date"`,
				`"advertise_orders"."advertise_package_id"`,
				`"a".*`,
			},
		},
		AdvertiseOrder: advertise_order.AdvertiseOrder{
			BusinessId: uid,
		},
	})

	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.BusinessAdvertiseOrderGetResponse_Data{
		Pagination: lib.Pagination(offset, limit, total),
		Result:     s.Model.ConvertAdvertiseOrderToProtos(res),
	}, nil
}

func (s *BusinessService) GetInvitationCode(ctx context.Context, req *pb.BusinessInvitationCodeGetRequest) (*pb.BusinessInvitationCodeGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetInvitationCode))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	bus, err := s.Model.GetBusinessById(ctx, req.XUserId)

	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	var code string
	if bus.InvitationCode == nil {
		code = ""
	} else {
		code = *bus.InvitationCode
	}

	return &pb.BusinessInvitationCodeGetResponse_Data{
		Code: code,
	}, nil
}

func (s *BusinessService) GetNumberFreeContact(ctx context.Context, req *pb.BusinessFreeContactGetRequest) (*pb.BusinessFreeContactGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetInvitationCode))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "Id"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	bus, err := s.Model.GetBusinessById(ctx, req.Id)

	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	var num int32
	if bus.FreeContact == nil {
		num = 0
	} else {
		num = *bus.FreeContact
	}

	return &pb.BusinessFreeContactGetResponse_Data{
		Number: num,
	}, nil
}

func (s BusinessService) BuyAdvertiseSetup(ctx context.Context, req *pb.BusinessBuyAdvertiseSetupPostRequest) (*pb.BusinessBuyAdvertiseSetupPostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.BuyAdvertiseSetup))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	p, err := s.Model.GetPaymentById(ctx, req.PaymentId)

	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	status, err := s.Model.BuyAdvertiseSetup(ctx, p.CustomerId, &req.Amount, p.PaymentMethodId)

	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.BusinessBuyAdvertiseSetupPostResponse_Data{
		Status: *status,
	}, nil
}

func (s BusinessService) BuyAdvertise(ctx context.Context, req *pb.BusinessBuyAdvertisePostRequest) (*pb.BusinessBuyAdvertisePostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.BuyAdvertise))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	ubid, err := lib.ToUUID(req.XUserId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	uapid, err := lib.ToUUID(req.AdvertisePackageId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	ucid, err := lib.ToUUID(req.CategoryId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	sr, err := s.Model.GetServiceForPayment(ctx, &service.Search{
		Service: service.Service{
			BusinessId: ubid,
			CategoryId: ucid,
		},
	})

	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	err = s.Model.CreateAdvertiseOrder(ctx, &advertise_transaction.AdvertiseTransaction{
		Name:         &req.PackageName,
		Price:        &req.Price,
		BannerUrl:    utils.SafeStrPtr(req.BannerUrl),
		Description:  &req.Description,
		Zipcode:      &req.Zipcode,
		CategoryName: &req.CategoryName,
	}, &advertise_order.AdvertiseOrder{
		BusinessId:         ubid,
		StartDate:          &req.StartDate,
		EndDate:            &req.EndDate,
		AdvertisePackageId: uapid,
		ServiceId:          sr.ID,
	})

	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.BusinessBuyAdvertisePostResponse_Data{}, nil
}

func (s BusinessService) ValidateBuyAdvertise(ctx context.Context, req *pb.BusinessValidateBuyAdvertisePostRequest) (*pb.BusinessValidateBuyAdvertisePostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ValidateBuyAdvertise))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId", "Zipcode", "CategoryId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	ucid, err := lib.ToUUID(req.CategoryId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	totalOrder, err := s.Model.GetTotalOrderForBuyValidate(ctx, &advertise_order.Search{
		CategoryId: ucid,
		Zipcode:    &req.Zipcode,
	})

	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	if *totalOrder > c.BUY_ADVERTISE_LIMIT-1 {
		err := xerrors.Errorf("%w", e.ExceedBuyAdvertiseLimit)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.BusinessValidateBuyAdvertisePostResponse_Data{}, nil
}

func (s *BusinessService) VerifyRefCode(ctx context.Context, req *pb.BusinessVerifyRefCodePutRequest) (*pb.BusinessVerifyRefCodePutResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.VerifyRefCode))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "Id", "Status"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	var statusVerify c.STATUS_VERIFY_REFERRAL_CODE
	if req.Status {
		statusVerify = c.STATUS_VERIFY_REFERRAL_CODE_ACCEPT
	} else {
		statusVerify = c.STATUS_VERIFY_REFERRAL_CODE_DENY
	}

	err := s.Model.UpdateBusiness(ctx, req.Id, &business.Business{
		RefStatus: utils.Int32Ptr(int32(statusVerify)),
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	if req.Status {
		bus, err := s.Model.GetBusinessById(ctx, req.Id)
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err, ctx)
			return nil, err
		}

		refBus, err := s.Model.GetBusinessByReferralCode(ctx, bus.RefCode)
		if err != nil || refBus == nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err, ctx)
			return nil, err
		}

		newFreeContact := utils.Int32Val(refBus.FreeContact) + 10

		err = s.Model.UpdateBusiness(ctx, refBus.ID, &business.Business{
			FreeContact: &newFreeContact,
		})
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err, ctx)
			return nil, err
		}
	}

	return &pb.BusinessVerifyRefCodePutResponse_Data{}, nil
}
