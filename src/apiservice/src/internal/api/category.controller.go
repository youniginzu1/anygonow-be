package api

import (
	"strconv"

	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/apiservice/src/pb"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var categorySet = wire.NewSet(wire.Struct(new(CategoryController), "*"), wire.Struct(new(CategoryService), "*"))

type CategoryController struct {
	S CategoryService
}

func (s *CategoryController) HandleGet(g *gin.Context) {
	q, err := strconv.Atoi(g.DefaultQuery("query", "0"))
	if err != nil {
		q = 0
	}
	req := pb.CategoriesGetRequest{
		Query:  c.QUERY_CATEGORY_ADMIN(q),
		Limit:  g.DefaultQuery("limit", "0"),
		Offset: g.DefaultQuery("offset", "0"),
		Name:   g.DefaultQuery("name", ""),
	}
	res, err := s.S.ListCategories(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *CategoryController) HandleByIdGet(g *gin.Context) {
	req := pb.CategoryGetRequest{
		Id: g.Param("id"),
	}
	res, err := s.S.GetCategory(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}
