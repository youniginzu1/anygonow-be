package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"reflect"
	"runtime"
	"strconv"

	"github.com/aqaurius6666/chatservice/src/internal/var/c"
	"github.com/aqaurius6666/chatservice/src/internal/var/e"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func RandomKey(name *string) *string {
	ext := path.Ext(*name)
	newKey := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	return &newKey
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
func ParseInt64Val(a string) int64 {
	if a == "" {
		return 0
	}
	v, err := strconv.ParseInt(a, 10, 64)
	if err != nil {
		return 0
	}
	return v
}
func ParseInt64Ptr(a string) *int64 {
	if a == "" {
		return nil
	}
	v, err := strconv.ParseInt(a, 10, 64)
	if err != nil {
		return nil
	}
	return &v
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

// func Pagination(offset, limit int32, total *int64) *pb.Pagination {
// 	if offset == 0 && limit == 0 && total == nil {
// 		return nil
// 	}
// 	return &pb.Pagination{
// 		Offset: offset,
// 		Limit:  limit,
// 		Total:  utils.Int64Val(total),
// 	}
// }

func GetFunctionName(funcname interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(funcname).Pointer()).Name()
}

func RecordError(span trace.Span, err error) {
	span.RecordError(fmt.Errorf("%+v", err))
	span.SetStatus(codes.Error, err.Error())
}
