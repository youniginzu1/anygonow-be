package middleware

import (
	"net/http"

	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/lib/unleash"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/apiservice/src/internal/var/e"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

type Auth interface {
	CheckAuth(g *gin.Context)
}

func (s *MiddlewareV1) CheckAuth(g *gin.Context) {
	if unleash.IsEnabled("apiservice.auth.by-pass") {
		s.byPassCheckAuth(g)
	} else {
		s.checkAuth(g)
	}

}

func (s *MiddlewareV1) checkAuth(g *gin.Context) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(lib.ParseGinContext(g), lib.GetFunctionName(s.checkAuth))
	defer span.End()

	var auth []byte
	var body []byte
	var method string
	var err error
	method = g.Request.Method
	auth = []byte(g.GetHeader("Authorization"))
	if len(auth) == 0 {
		lib.Unauthorized(g, xerrors.Errorf("%w", e.ErrMissingCertificate))
		return
	}
	if method != http.MethodGet {
		rawBody, ok := g.Get("rawBody")
		if !ok {
			lib.BadRequest(g, e.ErrExpectedBody)
			return
		}
		body = rawBody.([]byte)
	}
	id, role, err := s.Auth.CheckAuth(ctx, auth, body, method)
	if err != nil {
		lib.Unauthorized(g, xerrors.Errorf("%w", err))
		return
	}
	g.Set("userId", id)
	g.Set("role", role)
	g.Next()
}

func (s *MiddlewareV1) byPassCheckAuth(g *gin.Context) {
	id := g.Query("userId")
	role := g.Query("role")

	crole := c.ROLE(c.ROLE_value[role])
	g.Set("userId", id)
	g.Set("role", crole)
	g.Next()
}
