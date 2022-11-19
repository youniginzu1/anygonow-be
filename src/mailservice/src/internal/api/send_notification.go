package api

import (
	"context"

	"github.com/aqaurius6666/mailservice/src/internal/lib"
	"github.com/aqaurius6666/mailservice/src/internal/var/c"
	"github.com/aqaurius6666/mailservice/src/internal/var/e"
	"github.com/aqaurius6666/mailservice/src/pb/mailpb"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

func (s *ApiServer) SendNotification(ctx context.Context, req *mailpb.SendNotificationRequest) (*mailpb.SendNotificationResponse, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.SendNotification))
	defer span.End()

	if req.Body == "" || req.Message == "" || req.To == "" {
		err := xerrors.Errorf("%w", e.ErrInvalidRequest)
		lib.RecordError(span, err)
		panic(err)
	}
	return s.Model.SendNotification(ctx, req)
}
