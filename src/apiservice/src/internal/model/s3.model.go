package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

type S3Model interface {
	GetUploadUrl(ctx context.Context, key string, contentLength int64, userId string) (string, map[string]string, error)
}

func (s *ServerModel) GetUploadUrl(ctx context.Context, filename string, contentLength int64, userId string) (string, map[string]string, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetUploadUrl))
	defer span.End()

	contentType, err := lib.GetContentType(filename)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return "", nil, err
	}
	key := fmt.Sprintf("%s/%s", userId, filename)
	url, form, err := s.S3.GetPresignedPutObject(ctx, key, contentLength, contentType)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return "", nil, err
	}
	return strings.Replace(url, "http", "https", 1), form, nil
}
