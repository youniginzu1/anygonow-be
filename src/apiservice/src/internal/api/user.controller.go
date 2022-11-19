package api

import (
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/pb"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var userSet = wire.NewSet(wire.Struct(new(UserController), "*"), wire.Struct(new(UserService), "*"))

type UserController struct {
	S UserService
}

func (s UserController) HandleUserPost(g *gin.Context) {
	req := pb.UserPostRequest{}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	res, err := s.S.CreateUser(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s UserController) HandleGetById(g *gin.Context) {
	req := pb.UserGetRequest{
		Id: g.Param("id"),
	}
	res, err := s.S.GetById(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s UserController) HandlePut(g *gin.Context) {
	req := pb.UserPutRequest{}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.Id = g.Param("id")
	req.XUserId = g.GetString("userId")
	res, err := s.S.Update(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s UserController) HandleStateGet(g *gin.Context) {
	req := pb.UserStateGetRequest{}

	res, err := s.S.StateGet(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}
