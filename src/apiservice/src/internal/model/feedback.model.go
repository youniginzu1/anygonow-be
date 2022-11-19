package model

import (
	"context"

	"github.com/aqaurius6666/apiservice/src/internal/db/feedback"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/apiservice/src/pb"
	"github.com/aqaurius6666/go-utils/utils"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

var (
	_ FeedbackModel = (*ServerModel)(nil)
)

type FeedbackModel interface {
	GetRatingByBusiness(ctx context.Context, bid interface{}) (*feedback.Feedback, error)
	ListRatingByBusiness(ctx context.Context, bid interface{}) ([]*feedback.Feedback, error)

	ConvertRatingToProto(*feedback.Feedback) *pb.Rating
	ConvertRatingToProtos([]*feedback.Feedback) []*pb.Rating

	ListFeedbackByBusiness(context.Context, *feedback.Search) ([]*feedback.Feedback, error)
	TotalFeedbackByBusiness(context.Context, *feedback.Search) (*int64, error)

	ConvertFeedbackToProto(v *feedback.Feedback) *pb.Feedback
	ConvertFeedbackToProtos(v []*feedback.Feedback) []*pb.Feedback

	CreateFeedback(ctx context.Context, feedback *feedback.Feedback) (*feedback.Feedback, error)
	GetFeedback(ctx context.Context, search *feedback.Search) (*feedback.Feedback, error)
	UpdateFeedBack(ctx context.Context, search *feedback.Search, feedback *feedback.Feedback) error
}

func (s *ServerModel) ListRatingByBusiness(ctx context.Context, bid interface{}) ([]*feedback.Feedback, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListRatingByBusiness))
	defer span.End()

	uid, err := lib.ToUUID(bid)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	fb, err := s.Repo.ListRatingByBussiness(ctx, &feedback.Search{
		Feedback: feedback.Feedback{BusinessId: uid},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return fb, nil
}

func (s *ServerModel) GetRatingByBusiness(ctx context.Context, bid interface{}) (*feedback.Feedback, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetRatingByBusiness))
	defer span.End()

	uid, err := lib.ToUUID(bid)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	fb, err := s.Repo.SelectRatingByBussiness(ctx, &feedback.Search{
		Feedback: feedback.Feedback{BusinessId: uid},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return fb, nil
}

func (s *ServerModel) ListFeedbackByBusiness(ctx context.Context, search *feedback.Search) ([]*feedback.Feedback, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListFeedbackByBusiness))
	defer span.End()

	fb, err := s.Repo.ListFeedbacksByBusiness(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return fb, nil
}

func (s *ServerModel) TotalFeedbackByBusiness(ctx context.Context, search *feedback.Search) (*int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.TotalFeedbackByBusiness))
	defer span.End()

	c, err := s.Repo.TotalFeedbacksByBusiness(ctx, search)

	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return c, nil
}

func (s *ServerModel) ConvertRatingToProto(v *feedback.Feedback) *pb.Rating {
	upb := &pb.Rating{}
	if v.Rate != nil {
		upb.Rate = utils.Float32Val(v.Rate)
	}
	if v.Review != nil {
		upb.Review = *v.Review
	}
	if v.Request != nil {
		upb.Request = *v.Request
	}
	return upb
}

func (s *ServerModel) ConvertRatingToProtos(v []*feedback.Feedback) []*pb.Rating {
	arr := make([]*pb.Rating, 0)
	for _, vv := range v {
		arr = append(arr, s.ConvertRatingToProto(vv))
	}
	return arr
}

func (s *ServerModel) ConvertFeedbackToProto(v *feedback.Feedback) *pb.Feedback {
	upb := &pb.Feedback{}
	if v.Rate != nil {
		upb.Rate = utils.Float32Val(v.Rate)
	}
	if v.Comment != nil {
		upb.Comment = *v.Comment
	}
	if v.AvatarUrl != nil {
		upb.Image = *v.AvatarUrl
	}
	if v.ServiceName != nil {
		upb.ServiceOrder = *v.ServiceName
	}
	if v.CustomerName != nil {
		upb.CustomerName = *v.CustomerName
	}
	if v.CreatedAt != 0 {
		upb.CreatedAt = v.CreatedAt
	}
	return upb
}

func (s *ServerModel) ConvertFeedbackToProtos(v []*feedback.Feedback) []*pb.Feedback {
	arr := make([]*pb.Feedback, 0)
	for _, vv := range v {
		arr = append(arr, s.ConvertFeedbackToProto(vv))
	}
	return arr
}

func (s *ServerModel) CreateFeedback(ctx context.Context, feedback *feedback.Feedback) (*feedback.Feedback, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.TotalFeedbackByBusiness))
	defer span.End()

	fb, err := s.Repo.InsertFeedback(ctx, feedback)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return fb, nil
}

func (s *ServerModel) GetFeedback(ctx context.Context, search *feedback.Search) (*feedback.Feedback, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetFeedback))
	defer span.End()

	c, err := s.Repo.SelectFeedback(ctx, search)

	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return c, nil
}

func (s *ServerModel) UpdateFeedBack(ctx context.Context, search *feedback.Search, feedback *feedback.Feedback) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UpdateFeedBack))
	defer span.End()

	err := s.Repo.UpdateFeedback(ctx, search, feedback)

	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}

	return nil
}
