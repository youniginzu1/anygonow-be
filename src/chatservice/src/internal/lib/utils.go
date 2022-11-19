package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"

	"github.com/aqaurius6666/chatservice/src/internal/var/c"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SortBy []uuid.UUID

func (a SortBy) Len() int           { return len(a) }
func (a SortBy) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortBy) Less(i, j int) bool { return a[i].ID() < a[j].ID() }

type uuidWithPhoneNumber struct {
	id    uuid.UUID
	phone string
}

type sortableUUIDWithPhoneNumber []uuidWithPhoneNumber

func (a sortableUUIDWithPhoneNumber) Len() int           { return len(a) }
func (a sortableUUIDWithPhoneNumber) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortableUUIDWithPhoneNumber) Less(i, j int) bool { return a[i].id.ID() < a[j].id.ID() }
func StatusFromError(err error) error {
	err = xerrors.Errorf("%w", err)
	fmt.Printf("%+v", err)
	return status.Error(codes.Internal, err.Error())
}

func SortUUID(list []uuid.UUID) []uuid.UUID {
	sort.Sort(SortBy(list))
	return list
}

func IsSliceEqual(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func SortUUIDWithPhoneNumber(listUuid []uuid.UUID, listPhone []string) ([]uuid.UUID, []string) {
	list := make([]uuidWithPhoneNumber, 0)
	for i := range listUuid {
		list = append(list, uuidWithPhoneNumber{
			id:    listUuid[i],
			phone: listPhone[i],
		})
	}
	sort.Sort(sortableUUIDWithPhoneNumber(list))
	outUuid := make([]uuid.UUID, 0)
	outPhone := make([]string, 0)
	for _, uwpn := range list {
		outUuid = append(outUuid, uwpn.id)
		outPhone = append(outPhone, uwpn.phone)
	}
	return outUuid, outPhone
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

func StandardPhoneNumber(phone string) string {
	if regexp.MustCompile(`^\+1[0-9]{10}$`).MatchString(phone) {
		return phone
	}
	if regexp.MustCompile(`^[0-9]{10}$`).MatchString(phone) {
		return "+1" + phone
	}
	return fmt.Sprintf("+1%s", phone)
}
