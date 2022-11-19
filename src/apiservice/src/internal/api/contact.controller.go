package api

import (
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/pb"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var contactSet = wire.NewSet(wire.Struct(new(ContactController), "*"), wire.Struct(new(ContactService), "*"))

type ContactController struct {
	S ContactService
}

func (s *ContactController) HandleGetById(g *gin.Context) {
	req := pb.ContactGetRequest{
		Id: g.Param("id"),
	}
	res, err := s.S.GetById(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *ContactController) HandleListStates(g *gin.Context) {
	req := pb.StatesGetRequest{}
	res, err := s.S.ListStates(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *ContactController) HandlePut(g *gin.Context) {
	req := pb.ContactPutRequest{}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.XUserId = g.GetString("userId")
	req.Id = g.Param("id")
	res, err := s.S.Update(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *ContactController) HandleGetState(g *gin.Context) {
	req := pb.ContactGetRequest{
		Id: g.Param("id"),
	}
	res, err := s.S.GetStateById(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}