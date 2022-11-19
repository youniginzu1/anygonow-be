package api

import (
	"context"

	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/lib/validate"
	"github.com/aqaurius6666/apiservice/src/internal/model"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/apiservice/src/internal/var/e"
	"github.com/aqaurius6666/apiservice/src/pb"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

type ChatService struct {
	Model model.Server
}

func (s *ChatService) GetConversations(ctx context.Context, req *pb.ConversationPostRequest) (*pb.ConversationPostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetConversations))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	if len(req.ConversationIds) == 0 {
		return &pb.ConversationPostResponse_Data{}, nil
	}
	res, err := s.Model.GetConversation(ctx, req.XUserId, req.ConversationIds)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return res, nil
}
