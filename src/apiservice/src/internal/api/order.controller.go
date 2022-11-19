package api

import (
	"strconv"

	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/apiservice/src/pb"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var orderSet = wire.NewSet(wire.Struct(new(OrderController), "*"), wire.Struct(new(OrderService), "*"))

type OrderController struct {
	S OrderService
}

func (s OrderController) HandlePost(g *gin.Context) {
	req := pb.OrdersPostRequest{}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.XUserId = g.GetString("userId")
	res, err := s.S.CreateOrder(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *OrderController) HandleGet(g *gin.Context) {
	q, err := strconv.Atoi(g.DefaultQuery("status", "0"))
	if err != nil {
		q = 0
	}
	req := pb.OrdersGetRequest{
		XUserId:   g.GetString("userId"),
		Offset:    g.DefaultQuery("offset", "0"),
		Limit:     g.DefaultQuery("limit", "5"),
		Status:    c.ORDER_STATUS(q),
		ServiceId: g.DefaultQuery("serviceId", ""),
		Zipcode:   g.Query("zipcode"),
	}
	res, err := s.S.ListOrders(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *OrderController) HandleConnectPost(g *gin.Context) {
	req := pb.UpdateOrderStatusPostRequest{}

	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.XUserId = g.GetString("userId")
	req.XRole = lib.MustGetRole(g)

	res, err := s.S.UpdateOrderConnectedStatus(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *OrderController) HandleConnectAllPost(g *gin.Context) {
	req := pb.UpdateAllOrderStatusPostRequest{}

	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.XUserId = g.GetString("userId")
	req.XRole = lib.MustGetRole(g)

	res, err := s.S.UpdateAllOrderConnectedStatus(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *OrderController) HandleCompletePost(g *gin.Context) {
	req := pb.UpdateOrderStatusPostRequest{}

	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.XUserId = g.GetString("userId")
	req.XRole = lib.MustGetRole(g)

	res, err := s.S.UpdateOrderCompletedStatus(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *OrderController) HandleCancelPost(g *gin.Context) {
	req := pb.UpdateOrderStatusPostRequest{}

	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.XUserId = g.GetString("userId")
	req.XRole = lib.MustGetRole(g)

	res, err := s.S.UpdateOrderCancelledStatus(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *OrderController) HandleRejectPost(g *gin.Context) {
	req := pb.UpdateOrderStatusPostRequest{}

	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.XUserId = g.GetString("userId")
	req.XRole = lib.MustGetRole(g)

	res, err := s.S.UpdateOrderRejectedStatus(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *OrderController) HandleProjectsGet(g *gin.Context) {
	req := pb.UserProjectsGetRequest{
		Offset: g.DefaultQuery("offset", "0"),
		Limit:  g.DefaultQuery("limit", "15"),
	}
	req.XUserId = g.GetString("userId")

	res, err := s.S.GetProjects(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)

}

func (s *OrderController) HandleProjectCancelPost(g *gin.Context) {
	req := pb.CancelProjectPostRequest{}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.XUserId = g.GetString("userId")
	res, err := s.S.CancelProject(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)

}

func (s *OrderController) HandleBusinessesAlreadyOrderedGet(g *gin.Context) {
	req := pb.BusinessesAlreadyOrderedGetRequest{
		CategoryId: g.Query("categoryId"),
		Zipcode:    g.Query("zipcode"),
	}
	req.XUserId = g.GetString("userId")

	res, err := s.S.BusinessesAlreadyOrdered(g, &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}
