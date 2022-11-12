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

func (s *ApiServer) UnsubscribeNotification(ctx context.Context, req *mailpb.UnsubscribeNotificationRequest) (*mailpb.UnsubscribeNotificationResponse, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UnsubscribeNotification))
	defer span.End()
	if req.UserId == "" || req.DeviceId == "" {
		err := xerrors.Errorf("%w", e.ErrInvalidRequest)
		lib.RecordError(span, err)
		panic(err)
	}
	return s.Model.UnsubscribeNotification(ctx, req)
}
