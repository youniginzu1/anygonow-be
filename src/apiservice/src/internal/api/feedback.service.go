package api

import (
	"context"

	"github.com/aqaurius6666/apiservice/src/internal/db/feedback"
	"github.com/aqaurius6666/apiservice/src/internal/db/order"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/lib/validate"
	"github.com/aqaurius6666/apiservice/src/internal/model"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/apiservice/src/internal/var/e"
	"github.com/aqaurius6666/apiservice/src/pb"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/aqaurius6666/go-utils/utils"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

type FeedbackService struct {
	Model model.Server
}

func (s *FeedbackService) CreateFeedback(ctx context.Context, req *pb.FeedbacksPostRequest) (*pb.FeedbacksPostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CreateFeedback))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId", "OrderId", "Rate", "ServiceId", "BusinessId"); !ok {
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

	uid, err := lib.ToUUID(req.XUserId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	uoid, err := lib.ToUUID(req.OrderId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	usid, err := lib.ToUUID(req.ServiceId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	ubid, err := lib.ToUUID(req.BusinessId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	if *ord.Status != int32(c.ORDER_STATUS_COMPLETED) || ord.CustomerId != uid || ord.BusinessId != ubid || ord.ServiceId != usid {
		err = xerrors.Errorf("%w", e.ErrBodyInvalid)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	if check := true; *ord.IsReviewed == check {
		err = xerrors.Errorf("%w", e.ErrFeedbackExisted)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	err = s.Model.UpdateOrderById(ctx, ord.ID, &order.Order{
		IsReviewed: utils.BoolPtr(true),
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	fb, err := s.Model.CreateFeedback(ctx, &feedback.Feedback{
		UserId:     uid,
		OrderId:    uoid,
		Rate:       &req.Rate,
		Comment:    &req.Comment,
		ServiceId:  usid,
		BusinessId: ubid,
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &pb.FeedbacksPostResponse_Data{
		Feedback: s.Model.ConvertFeedbackToProto(fb),
	}, nil
}

func (s *FeedbackService) GetFeedBack(ctx context.Context, req *pb.FeedbackGetRequest) (*pb.FeedbackGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetFeedBack))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId", "OrderId"); !ok {
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

	uoid, err := lib.ToUUID(req.OrderId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	fb, err := s.Model.GetFeedback(ctx, &feedback.Search{
		Feedback: feedback.Feedback{
			OrderId: uoid,
			UserId:  uid,
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.FeedbackGetResponse_Data{
		Feedback: s.Model.ConvertFeedbackToProto(fb),
	}, nil
}

func (s *FeedbackService) UpdateFeedBack(ctx context.Context, req *pb.FeedbackPutRequest) (*pb.FeedbackPutResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UpdateFeedBack))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId", "Rate", "Comment", "Id"); !ok {
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

	fid, err := lib.ToUUID(req.Id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	_, err = s.Model.GetFeedback(ctx, &feedback.Search{
		Feedback: feedback.Feedback{
			BaseModel: database.BaseModel{
				ID: fid,
			},
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	err = s.Model.UpdateFeedBack(ctx, &feedback.Search{
		Feedback: feedback.Feedback{
			BaseModel: database.BaseModel{
				ID: fid,
			},
			UserId: uid,
		},
	}, &feedback.Feedback{
		Comment: utils.SafeStrPtr(req.Comment),
		Rate:    utils.Float32Ptr(req.Rate),
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.FeedbackPutResponse_Data{}, nil
}
