package middleware

import (
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/apiservice/src/internal/var/e"
	"github.com/gin-gonic/gin"
)

type Role interface {
	OnlyAdmin() gin.HandlerFunc
	Only(...c.ROLE) gin.HandlerFunc
}

func (l *MiddlewareV1) OnlyAdmin() gin.HandlerFunc {
	return func(g *gin.Context) {
		r := lib.MustGetRole(g)
		if r != c.ROLE_ADMIN {
			lib.Unauthorized(g, e.ErrNoPermission)
			return
		}
		g.Next()
	}
}
func (l *MiddlewareV1) Only(ro ...c.ROLE) gin.HandlerFunc {
	return func(g *gin.Context) {
		r := lib.MustGetRole(g)
		shouldPass := false
		for _, ros := range ro {
			if r == ros {
				shouldPass = true
			}
		}
		if !shouldPass {
			lib.Unauthorized(g, e.ErrNoPermission)
		}
		g.Next()
	}
}
