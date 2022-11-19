package middleware

import (
	"encoding/json"
	"time"

	"github.com/aqaurius6666/chatservice/src/internal/lib/unleash"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Logger interface {
	LoggerMid() gin.HandlerFunc
}

func (l *MiddlewareV1) LoggerMid() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		latencyTime := time.Since(startTime)
		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		statusCode := c.Writer.Status()
		reqLogger := l.Logger.WithFields(logrus.Fields{
			"status":  statusCode,
			"latency": latencyTime,
			"method":  reqMethod,
			"path":    reqUri,
		})

		if body, ok := c.Get("body"); ok && unleash.IsEnabled("apiservice.debug.log-body") {
			reqLogger = reqLogger.WithField("body", string(body.(json.RawMessage)))
		}
		if err, ok := c.Get("error"); ok && unleash.IsEnabled("apiservice.debug.log-error") {
			reqLogger.Errorf("%+v", err)
		} else if unleash.IsEnabled("apiservice.debug.log-request") {
			reqLogger.Info()
		}
	}
}
