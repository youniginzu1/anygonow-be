package lib

import (
	"fmt"
	"reflect"
	"runtime"
	"time"

	"github.com/aqaurius6666/cronjob/src/internal/var/e"
	"github.com/google/uuid"
	otelcodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/xerrors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToUUID(id interface{}) (uuid.UUID, error) {
	switch tmp := id.(type) {
	case uuid.UUID:
		return tmp, nil
	case string:
		tmp2, err := uuid.Parse(tmp)
		if err != nil {
			return uuid.Nil, e.ErrIdInvalidFormat
		}
		return tmp2, nil
	default:
		return uuid.Nil, e.ErrIdInvalidFormat
	}
}

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

func GetOneMonthRange(n int64) (prev, now int64) {
	lastMonth := n - int64(28*24*time.Hour.Milliseconds())
	return lastMonth, n
}
func SafeBoolPtr(a bool) *bool {
	if a {
		return nil
	}
	return &a
}
func UsdFloatToCent(in *float32) *int64 {
	a := *in * 100
	in64 := int64(a)
	return &in64
}

func GetTimeRange(quantity int, unit string) time.Duration {
	if unit == "HOUR" {
		return time.Duration(quantity) * time.Hour
	}
	if unit == "DAY" {
		return time.Hour * time.Duration(quantity) * 24
	}
	if unit == "WEEK" {
		return 7 * 24 * time.Hour * time.Duration(quantity)
	}
	if unit == "MONTH" {
		return 28 * 24 * time.Hour * time.Duration(quantity)
	}
	return 0
}
