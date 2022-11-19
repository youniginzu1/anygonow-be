package middleware

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/e"
	"github.com/gin-gonic/gin"
)

type Body interface {
	HandleBody() gin.HandlerFunc
}

func (l *MiddlewareV1) HandleBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}
		var body struct {
			Data json.RawMessage `json:"data"`
		}
		bz, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			lib.BadRequest(c, e.ErrExpectedBody)
			return
		}
		if err := json.Unmarshal(bz, &body); err != nil {
			lib.BadRequest(c, e.ErrExpectedBody)
			return
		}
		c.Set("rawBody", bz)
		c.Set("body", body.Data)
		c.Next()

	}
}
