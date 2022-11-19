package lib

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"path"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/apiservice/src/internal/var/e"
	"github.com/aqaurius6666/apiservice/src/pb"
	"github.com/aqaurius6666/go-utils/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/xerrors"
)

func SetBody(c *gin.Context) error {
	if c.Request.Method == http.MethodGet {
		return nil
	}
	buf := &bytes.Buffer{}

	_, err := buf.ReadFrom(c.Request.Body)
	if err != nil {
		return xerrors.Errorf("%w", e.ErrMissingBody)
	}
	c.Set("body", buf.Bytes())
	c.Request.Body = ioutil.NopCloser(buf)

	return nil
}

func GetBody(g *gin.Context, target interface{}) error {
	var bbody []byte
	var err error
	ibody, ok := g.Get("body")
	if !ok {
		return xerrors.Errorf("%w", e.ErrMissingBody)
	} else {
		bbody = ibody.(json.RawMessage)
	}

	err = json.Unmarshal(bbody, target)
	if err != nil {
		return xerrors.Errorf("%w", e.ErrMissingBody)
	}
	return nil

}

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

func MustGetRole(g *gin.Context) c.ROLE {
	r, ok := g.Get("role")
	if !ok {
		panic("not role")
	}
	ro, ok := r.(c.ROLE)
	if !ok {
		panic("not role")
	}
	return ro
}

func GetContentType(name string) (string, error) {
	ext := path.Ext(name)
	typ := mime.TypeByExtension(ext)

	if !strings.Contains(typ, "image") {
		return "", e.ErrInvalidFile
	}
	return typ, nil
}

func RandomKey(name string) string {
	ext := path.Ext(name)
	newKey := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	return newKey
}

func ParseInt32Val(a string) int32 {
	if a == "" {
		return 0
	}
	v, err := strconv.Atoi(a)
	if err != nil {
		return 0
	}
	i32 := int32(v)
	return i32
}

func ParseInt32Ptr(a string) *int32 {
	if a == "" {
		return nil
	}
	v, err := strconv.Atoi(a)
	if err != nil {
		return nil
	}
	i32 := int32(v)
	return &i32
}

func ParseIntPtr(a string) *int {
	if a == "" {
		return nil
	}
	v, err := strconv.Atoi(a)
	if err != nil {
		return nil
	}
	return &v
}

func ParseUUID(a string) uuid.UUID {
	uid, err := uuid.Parse(a)
	if err != nil {
		return uuid.Nil
	}
	return uid
}

func Pagination(offset, limit int32, total *int64) *pb.Pagination {
	if offset == 0 && limit == 0 && total == nil {
		return nil
	}
	return &pb.Pagination{
		Offset: offset,
		Limit:  limit,
		Total:  utils.Int64Val(total),
	}
}

func GetFunctionName(funcname interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(funcname).Pointer()).Name()
}

func RecordError(span trace.Span, err error, ctxs ...context.Context) {
	span.RecordError(fmt.Errorf("%+v", err))
	span.SetStatus(codes.Error, err.Error())
	for _, ctx := range ctxs {
		if body := ctx.Value(GinKey("body")); body != nil {
			span.SetAttributes(attribute.KeyValue{
				Key:   attribute.Key("request.body"),
				Value: attribute.StringValue(string(body.(json.RawMessage))),
			})
		}
	}
}

type GinKey string

func ParseGinContext(c *gin.Context) context.Context {
	ctx := c.Request.Context()
	for k, v := range c.Keys {
		ctx = context.WithValue(ctx, GinKey(k), v)
	}
	return ctx
}

var (
	aDay  = time.Hour * 24
	aWeek = aDay * 7
)

func GetTimeRange(n int32) (left int64, right int64) {
	now := time.Now()
	weekday := int64(now.Weekday())
	left = now.UnixMilli() - now.UnixMilli()%aDay.Milliseconds() - weekday*aDay.Milliseconds()
	if n == 0 {
		return left, now.UnixMilli()
	}
	if n == 5 {
		return 0, now.UnixMilli()
	}
	left -= int64(n) * aWeek.Milliseconds()
	return left, left + aWeek.Milliseconds()
}

func ParseInt64Ptr(a string) *int64 {
	if a == "" {
		return nil
	}
	v, err := strconv.Atoi(a)
	if err != nil {
		return nil
	}
	i64 := int64(v)
	return &i64
}

func ParseInt64Val(a string) int64 {
	if a == "" {
		return 0
	}
	v, err := strconv.Atoi(a)
	if err != nil {
		return 0
	}
	i64 := int64(v)
	return i64
}

func ExtractUrlBeforeInsert(url string) *string {
	if url == "" {
		return nil
	}
	if !strings.Contains(url, BaseUrl) {
		return nil
	}
	regex := regexp.MustCompile(BaseUrl)
	extractedUrl := regex.ReplaceAllLiteralString(url, "")

	return utils.SafeStrPtr(extractedUrl)
}

func StandardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func UsdFloatToCent(in *float32) *int64 {
	a := *in * 100
	in64 := int64(a)
	return &in64
}

func CentToUsd(in *int64) *float32 {
	a := float32(*in) / 100
	return &a
}

func GenerateCellExcelData(a string, i int) string {
	return a + strconv.Itoa(i+EXCEL_ROW_INDEX_DATA_START)
}

func GenerateCellExcelHeader(a string) string {
	return a + strconv.Itoa(EXCEL_ROW_INDEX_HEADER)
}

func FormatMillisecondsToDate(miliSec int64) string {
	s := time.UnixMilli(miliSec)
	return s.Format("2006-01-02 15:04:05")
}

func Float64PtrToInt64Cent(a *float64) int64 {
	if a == nil {
		return 0
	}
	return int64(utils.Float64Val(a) * 100)
}
