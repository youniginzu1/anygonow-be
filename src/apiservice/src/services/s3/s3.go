package s3

import (
	"context"
	"io"
	"time"

	"github.com/google/wire"
)

type BucketName string

var (
	// Set = wire.NewSet(wire.Bind(new(Service), new(V1)), ConnectS3, NewV2, ConnectSigner, LoadConfig)
	Set = wire.NewSet(wire.Bind(new(Service), new(*V2)), NewV2)
)

var (
	timeout = 3 * time.Second
)

type Service interface {
	GetObject(key *string) (io.ReadCloser, error)
	GetPresignedPutObject(ctx context.Context, key string, contentLength int64, contentTYpe string) (url string, form map[string]string, err error)
}
