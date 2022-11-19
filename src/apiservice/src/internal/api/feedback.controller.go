package api

import (
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/pb"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var feedbackSet = wire.NewSet(wire.Struct(new(FeedbackController), "*"), wire.Struct(new(FeedbackService), "*"))

type FeedbackController struct {
	S FeedbackService
}

func (s *FeedbackController) HandlePost(g *gin.Context) {
	req := pb.FeedbacksPostRequest{}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.XUserId = g.GetString("userId")

	res, err := s.S.CreateFeedback(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *FeedbackController) HandleGet(g *gin.Context) {
	req := pb.FeedbackGetRequest{
		OrderId: g.Query("orderId"),
	}
	req.XUserId = g.GetString("userId")

	res, err := s.S.GetFeedBack(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *FeedbackController) HandlePut(g *gin.Context) {
	req := pb.FeedbackPutRequest{
		Id: g.Param("id"),
	}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.XUserId = g.GetString("userId")

	res, err := s.S.UpdateFeedBack(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}
