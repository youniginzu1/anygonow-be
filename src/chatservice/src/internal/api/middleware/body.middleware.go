package middleware

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"

	"github.com/aqaurius6666/chatservice/src/internal/lib"
	"github.com/aqaurius6666/chatservice/src/internal/var/e"
	"github.com/aqaurius6666/chatservice/src/services/twilloclient"
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
		contentType := c.GetHeader("content-type")

		switch contentType {

		case "application/x-www-form-urlencoded":
			// bz, err := ioutil.ReadAll(c.Request.Body)
			// if err != nil {
			// 	lib.BadRequest(c, e.ErrExpectedBody)
			// 	return
			// }
			var messageCallback twilloclient.TwilloMessageCallbackData
			err := c.Bind(&messageCallback)
			if err != nil {
				lib.BadRequest(c, e.ErrExpectedBody)
				return
			}
			bz, err := xml.Marshal(messageCallback)
			if err != nil {
				lib.BadRequest(c, e.ErrExpectedBody)
				return
			}
			c.Set("twillo-message", messageCallback)
			c.Set("raw-body", bz)

		default:
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
		}

		c.Next()

	}
}
