package api

import (
	"context"
	"fmt"

	"github.com/aqaurius6666/apiservice/src/internal/db/advertise_order"
	"github.com/aqaurius6666/apiservice/src/internal/db/business"
	"github.com/aqaurius6666/apiservice/src/internal/db/category"
	"github.com/aqaurius6666/apiservice/src/internal/db/transaction"
	"github.com/aqaurius6666/apiservice/src/internal/db/user"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/lib/validate"
	"github.com/aqaurius6666/apiservice/src/internal/model"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/apiservice/src/internal/var/e"
	"github.com/aqaurius6666/apiservice/src/pb"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/aqaurius6666/go-utils/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

var FREE_CONTACT_MONEY = 5

type IndexService struct {
	Model  model.Server
	Logger *logrus.Logger
}

func (s *IndexService) GetUploadUrl(ctx context.Context, req *pb.UploadUrlPostRequest) (interface{}, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetUploadUrl))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId", "Filename", "ContentLength"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	url, form, err := s.Model.GetUploadUrl(ctx, req.Filename, req.ContentLength, req.XUserId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return gin.H{
		"url":  url,
		"form": form,
	}, nil
}

func (s *IndexService) GetStripeKey(ctx context.Context, req *pb.StripeKeyGetRequest) (*pb.StripeKeyGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetStripeKey))
	defer span.End()

	key, err := s.Model.GetStripeKey(ctx)
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}

	return &pb.StripeKeyGetResponse_Data{
		Key: *key,
	}, nil
}
func (s *IndexService) SubscribeNotification(ctx context.Context, req *pb.SubscribePostRequest) (*pb.SubscribePostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.SubscribeNotification))
	defer span.End()
	if f, ok := validate.RequiredFields(req, "XUserId", "DeviceId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	err := s.Model.UnsubscribeNotification(ctx, req.XUserId, req.DeviceId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	return &pb.SubscribePostResponse_Data{}, nil
}

func (s *IndexService) UnsubscribeNotification(ctx context.Context, req *pb.UnsubscribePostRequest) (*pb.UnsubscribePostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UnsubscribeNotification))
	defer span.End()
	if f, ok := validate.RequiredFields(req, "XUserId", "DeviceId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	err := s.Model.SubscribeNotification(ctx, req.XUserId, req.DeviceId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return nil, err
	}
	return &pb.UnsubscribePostResponse_Data{}, nil
}

func (s *IndexService) GetPaymentMethodInfo(ctx context.Context, req *pb.StripePaymentMethodGetRequest) (interface{}, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetPaymentMethodInfo))
	defer span.End()
	if f, ok := validate.RequiredFields(req, "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	pay, err := s.Model.GetPaymentByBusinessId(ctx, req.XUserId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	if utils.SafeStrPtr(utils.StrVal(pay.PaymentMethodId)) == nil {
		err = xerrors.Errorf("%w", e.ErrPayment)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	pm, err := s.Model.GetPaymentMethodInfo(ctx, pay.PaymentMethodId)
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}

	return pm, nil
}

func (s *IndexService) UpdateTransactionWebHook(ctx context.Context, paymentIntent stripe.PaymentIntent) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UpdateTransactionWebHook))
	defer span.End()

	if err := s.Model.UpdateTransaction(ctx, &transaction.Search{
		Transaction: transaction.Transaction{
			PaymentIntentId: utils.SafeStrPtr(paymentIntent.ID),
			IsPaid:          utils.SafeBoolPtr(false),
			IsFree:          utils.SafeBoolPtr(false),
		},
	}, &transaction.Transaction{
		IsPaid: utils.BoolPtr(true),
	}); err != nil {
		err = xerrors.Errorf("%w", e.ErrPayment)
		lib.RecordError(span, err, ctx)
		return err
	}

	trans, err := s.Model.GetTransactions(ctx, &transaction.Search{
		DefaultSearchModel: database.DefaultSearchModel{
			Fields: []string{
				`"transactions"."id"`,
				`"transactions"."updated_at"`,
				`"transactions"."created_at"`,
				`"transactions"."deleted_at"`,
				`"transactions"."order_id"`,
				`"transactions"."business_id"`,
				`"transactions"."fee"`,
				`"transactions"."is_free"`,
				`"transactions"."payment_intent_id"`,
			},
		},
		Transaction: transaction.Transaction{
			PaymentIntentId: &paymentIntent.ID,
		},
		Query: utils.Int32Ptr(int32(c.SORT_TRANSACTION_ALL_TIME)),
	})
	if err != nil {
		err = xerrors.Errorf("%w", e.ErrPayment)
		lib.RecordError(span, err, ctx)
		return err
	}

	var bId uuid.UUID
	if len(trans) > 0 {
		bId = trans[0].BusinessId
	} else {
		pay, err := s.Model.GetPaymentByCustomerId(ctx, &paymentIntent.Customer.ID)
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err, ctx)
			return err
		}
		bId = pay.BusinessId
	}

	if bId == uuid.Nil {
		err = xerrors.Errorf("%w", e.ErrPayment)
		lib.RecordError(span, err, ctx)
		return err
	}

	notifyFunc := func(handymanId string, paymentIntentId string) {
		ctx := context.TODO()
		ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName("notifyFunc"))
		defer span.End()

		amount, err := s.Model.GetPaymentIntentAmount(ctx, &paymentIntentId)
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err)
			return
		}
		err = s.Model.SendFeeNotification(ctx, handymanId, float32(amount)/100)
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err)
			return
		}
	}

	go notifyFunc(bId.String(), paymentIntent.ID)

	// Get total fee of businesses have just paid
	fee1, err := s.Model.TotalFee(ctx, &transaction.Search{
		Transaction: transaction.Transaction{
			BusinessId: bId,
		},
		Query: utils.Int32Ptr(int32(c.SORT_TRANSACTION_ALL_TIME)),
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}

	fee2, err := s.Model.TotalFeeAdvertise(ctx, &advertise_order.Search{
		AdvertiseOrder: advertise_order.AdvertiseOrder{
			BusinessId: bId,
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}

	// Get business info have just paid
	bus, err := s.Model.GetBusinessById(ctx, bId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}

	if utils.Int64Val(fee1)+lib.Float64PtrToInt64Cent(fee2) >= int64(FREE_CONTACT_MONEY) && utils.Int32Val(bus.RefStatus) == int32(c.STATUS_VERIFY_REFERRAL_CODE_VERIFYING) {
		// Get business info which invited the business above
		refBus, err := s.Model.GetBusinessByReferralCode(ctx, bus.RefCode)
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err, ctx)
			return err
		}

		// Update that the business's refcode have just paid is accepted
		err = s.Model.UpdateBusiness(ctx, bus.ID, &business.Business{
			RefStatus: utils.Int32Ptr(int32(c.STATUS_VERIFY_REFERRAL_CODE_ACCEPT)),
		})
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err, ctx)
			return err
		}

		// Add free contact
		err = s.Model.UpdateBusiness(ctx, refBus.ID, &business.Business{
			FreeContact: utils.Int32Ptr(utils.Int32Val(refBus.FreeContact) + 10),
		})
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err, ctx)
			return err
		}
	}

	return nil
}
func (s *IndexService) GetHomepageStatistic(ctx context.Context, req *pb.StatisticGetRequest) (*pb.StatisticGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetPaymentMethodInfo))
	defer span.End()
	getTotalBusiness := func(ctx context.Context, search *business.Search) <-chan *int64 {
		ch := make(chan *int64)
		go func() {
			totalBusinesses, err := s.Model.TotalBusinesses(ctx, search)
			if err != nil {
				err = xerrors.Errorf("%w", err)
				lib.RecordError(span, err, ctx)
				ch <- utils.Int64Ptr(0)
			}
			ch <- totalBusinesses
		}()
		return ch
	}(ctx, &business.Search{})
	getTotalUser := func(ctx context.Context, search *user.Search) <-chan *int64 {
		ch := make(chan *int64)
		go func() {
			totalCustomers, err := s.Model.TotalUsers(ctx, &user.Search{})
			if err != nil {
				err = xerrors.Errorf("%w", err)
				lib.RecordError(span, err, ctx)
				ch <- utils.Int64Ptr(0)
			}
			ch <- totalCustomers
		}()
		return ch
	}(ctx, &user.Search{})

	getTotalCategories := func(ctx context.Context, search *category.Search) <-chan *int64 {
		ch := make(chan *int64)
		go func() {
			totalCategories, err := s.Model.TotalCategory(ctx, &category.Search{})
			if err != nil {
				err = xerrors.Errorf("%w", err)
				lib.RecordError(span, err, ctx)
				ch <- utils.Int64Ptr(0)
			}
			ch <- totalCategories
		}()
		return ch
	}(ctx, &category.Search{})

	getTotalCities := func(ctx context.Context, search *business.Search) <-chan *int64 {
		ch := make(chan *int64)
		go func() {
			totalCities, err := s.Model.GetTotalZipcodes(ctx, &business.Search{})
			if err != nil {
				err = xerrors.Errorf("%w", err)
				lib.RecordError(span, err, ctx)
				ch <- utils.Int64Ptr(0)
			}
			ch <- totalCities
		}()
		return ch
	}(ctx, &business.Search{})

	totalBusinesses, totalCustomers, totalCategories, totalCities := <-getTotalBusiness, <-getTotalUser, <-getTotalCategories, <-getTotalCities

	return &pb.StatisticGetResponse_Data{
		BusinessQuantity: *totalBusinesses,
		CustomerQuantity: *totalCustomers,
		ServiceQuantity:  *totalCategories,
		CityQuantity:     *totalCities,
	}, nil
}

func (s *IndexService) ValidateMail(ctx context.Context, req *pb.ValidateMailGetRequest) (*pb.ValidateMailGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ValidateMail))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "Mail"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	_, err1 := s.Model.GetUser(ctx, &user.Search{
		Query: &req.Mail,
	})

	_, err2 := s.Model.GetBusiness(ctx, &business.Search{
		ValidateMail: &req.Mail,
		DefaultSearchModel: database.DefaultSearchModel{
			Fields: []string{`"businesses"."id"`},
		},
	})
	if err1 != nil && err2 != nil {
		fmt.Println(err1, err2)
		return &pb.ValidateMailGetResponse_Data{
			IsValidate: true,
		}, nil
	}

	return &pb.ValidateMailGetResponse_Data{
		IsValidate: false,
	}, nil
}
