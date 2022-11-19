package api

import (
	"github.com/aqaurius6666/chatservice/src/internal/lib"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var indexSet = wire.NewSet(wire.Struct(new(IndexController), "*"))

type IndexController struct {
}

func (s *IndexController) HandleIndexGet(g *gin.Context) {
	lib.Success(g, "Go go bruh bruh ...")
}
