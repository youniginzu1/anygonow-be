package model

import (
	"context"

	"github.com/aqaurius6666/apiservice/src/internal/db/payment"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/apiservice/src/internal/var/e"
	"github.com/aqaurius6666/apiservice/src/pb"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

var (
	_ PaymentModel = (*ServerModel)(nil)
)

type PaymentModel interface {
	GetStripeKey(context.Context) (*string, error)
	GetPaymentMethodInfo(context.Context, *string) (*pb.PaymentMethodInfo, error)
	NewCustomer(ctx context.Context, id string, mail string) (*payment.Payment, error)
	GetPaymentByBusinessId(context.Context, interface{}) (*payment.Payment, error)
	SetupPaymentMethod(context.Context, *string) (*string, error)
	UpdatePayment(ctx context.Context, search *payment.Search, value *payment.Payment) error
	ConvertPaymentToProto(p *payment.Payment) *pb.Payment
	ConvertPaymentProtos(v []*payment.Payment) []*pb.Payment
	BuyAdvertiseSetup(ctx context.Context, customerId *string, amount *int64, paymentMethodId *string) (*string, error)
	GetPaymentById(context.Context, interface{}) (*payment.Payment, error)
	GetPaymentByCustomerId(context.Context, *string) (*payment.Payment, error)
	GetPaymentIntentAmount(context.Context, *string) (int64, error)
}

func (s *ServerModel) GetPaymentIntentAmount(ctx context.Context, id *string) (int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetPaymentIntentAmount))
	defer span.End()

	amount, err := s.Payment.GetPaymentIntentAmount(ctx, id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		return 0, err
	}
	return amount, nil
}

func (s *ServerModel) NewCustomer(ctx context.Context, id string, mail string) (*payment.Payment, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.NewCustomer))
	defer span.End()

	key, err := s.Payment.CreateCustomer(ctx, id, mail)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	pay, err := s.Repo.InsertPayment(ctx, &payment.Payment{
		BusinessId: uuid.MustParse(id),
		CustomerId: &key,
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return pay, err
}

func (s *ServerModel) GetPaymentMethodInfo(ctx context.Context, id *string) (*pb.PaymentMethodInfo, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.NewCustomer))
	defer span.End()

	pm, err := s.Payment.GetPaymentMethodInfo(ctx, id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return pm, err
}
func (s *ServerModel) SetupPaymentMethod(ctx context.Context, id *string) (*string, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.SetupPaymentMethod))
	defer span.End()
	if id == nil {
		return nil, e.ErrIdInvalidFormat
	}

	cs, err := s.Payment.SetupIntent(ctx, id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return cs, nil
}

func (s *ServerModel) GetStripeKey(ctx context.Context) (*string, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetStripeKey))
	defer span.End()

	key, err := s.Payment.GetPublicKey(ctx)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &key, err
}

func (s *ServerModel) GetPaymentByBusinessId(ctx context.Context, id interface{}) (*payment.Payment, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetPaymentByBusinessId))
	defer span.End()

	bid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	pay, err := s.Repo.SelectPayment(ctx, &payment.Search{
		Payment: payment.Payment{
			BusinessId: bid,
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return pay, err
}

func (s *ServerModel) SetupPayment(ctx context.Context, cusId *string) (*string, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.SetupPayment))
	defer span.End()

	secretClient, err := s.Payment.SetupIntent(ctx, cusId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return secretClient, err
}

func (s *ServerModel) UpdatePayment(ctx context.Context, search *payment.Search, value *payment.Payment) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UpdatePayment))
	defer span.End()

	err := s.Repo.UpdatePayment(ctx, search, value)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}

func (s *ServerModel) ConvertPaymentToProto(p *payment.Payment) *pb.Payment {
	ppb := &pb.Payment{}
	if p.PaymentMethodId != nil {
		ppb.PaymentMethodId = *p.PaymentMethodId
	}
	if p.BaseModel.ID != uuid.Nil {
		ppb.Id = p.BaseModel.ID.String()
	}
	return ppb
}

func (s *ServerModel) ConvertPaymentProtos(v []*payment.Payment) []*pb.Payment {
	arr := make([]*pb.Payment, 0)
	for _, vv := range v {
		arr = append(arr, s.ConvertPaymentToProto(vv))
	}
	return arr
}

func (s *ServerModel) BuyAdvertiseSetup(ctx context.Context, customerId *string, amount *int64, paymentMethodId *string) (*string, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.BuyAdvertiseSetup))
	defer span.End()

	cs, err := s.Payment.PaymentIntent(ctx, customerId, amount, paymentMethodId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return cs, nil
}

func (s *ServerModel) GetPaymentById(ctx context.Context, id interface{}) (*payment.Payment, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetPaymentById))
	defer span.End()

	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	pay, err := s.Repo.SelectPayment(ctx, &payment.Search{
		Payment: payment.Payment{
			BaseModel: database.BaseModel{
				ID: uid,
			},
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return pay, err
}

func (s *ServerModel) GetPaymentByCustomerId(ctx context.Context, cusId *string) (*payment.Payment, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetPaymentByCustomerId))
	defer span.End()

	pay, err := s.Repo.SelectPayment(ctx, &payment.Search{
		Payment: payment.Payment{
			CustomerId: cusId,
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return pay, err
}
