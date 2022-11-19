package s3

import (
	"context"
	"io"
	"time"

	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

type V2 struct {
	Ctx    context.Context
	Client *minio.Client
	Bucket BucketName
}

func NewV2(ctx context.Context, bucket BucketName) (*V2, error) {
	client, err := minio.New("s3.amazonaws.com", &minio.Options{
		Creds: credentials.NewEnvAWS(),
	})
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}

	return &V2{
		Ctx:    ctx,
		Client: client,
		Bucket: bucket,
	}, nil
}
func (s *V2) GetObject(key *string) (io.ReadCloser, error) {
	panic("not implemented")
}

func (s *V2) GetPresignedPutObject(ctx context.Context, key string, contentLength int64, contentType string) (url string, form map[string]string, err error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetPresignedPutObject))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	policy := minio.NewPostPolicy()
	policy.SetBucket(string(s.Bucket))
	policy.SetExpires(time.Now().Add(c.PRESIGNED_URL_EXPIRE_TIME))
	policy.SetKey(key)
	policy.SetSuccessStatusAction("200")
	policy.SetContentType(contentType)
	policy.SetContentLengthRange(1, contentLength)

	u, form, err := s.Client.PresignedPostPolicy(ctx, policy)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return "", nil, err
	}

	return u.String(), form, nil
}
