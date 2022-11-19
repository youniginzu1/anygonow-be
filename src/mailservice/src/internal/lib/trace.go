package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"

	"github.com/aqaurius6666/mailservice/src/internal/var/c"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func StatusFromError(err error) error {
	err = xerrors.Errorf("%w", err)
	fmt.Printf("%+v", err)
	return status.Error(codes.Internal, err.Error())
}
func GetFunctionName(funcname interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(funcname).Pointer()).Name()
}

func RecordError(span trace.Span, err error) {
	span.RecordError(fmt.Errorf("%+v", err))
	span.SetStatus(otelcodes.Error, err.Error())
}

func UnaryServerLogRequestInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, GetFunctionName(UnaryServerLogRequestInterceptor))
	defer span.End()
	bodyStr, err := json.Marshal(req)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		RecordError(span, err)
		panic(err)
	}
	span.SetAttributes(attribute.KeyValue{
		Key:   attribute.Key("request.body"),
		Value: attribute.StringValue(string(bodyStr)),
	})
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	resp, err = handler(ctx, req)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		RecordError(span, err)
		panic(err)
	}
	return resp, err
}
