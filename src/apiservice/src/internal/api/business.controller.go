package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/apiservice/src/pb"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var businessSet = wire.NewSet(wire.Struct(new(BusinessController), "*"), wire.Struct(new(BusinessService), "*"))

type BusinessController struct {
	S BusinessService
}

func (s *BusinessController) HandlePost(g *gin.Context) {
	req := pb.BusinessPostRequest{}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	res, err := s.S.CreateBusiness(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *BusinessController) HandleRatingGet(g *gin.Context) {
	req := pb.BusinessRatingGetRequest{
		Id: g.Param("id"),
	}

	res, err := s.S.GetRatings(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *BusinessController) HandleFeedbackGet(g *gin.Context) {
	req := pb.BusinessFeedbacksGetRequest{
		Id:     g.Param("id"),
		Offset: g.DefaultQuery("offset", "0"),
		Limit:  g.DefaultQuery("limit", "5"),
	}
	res, err := s.S.GetFeedbacks(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *BusinessController) HandleGetById(g *gin.Context) {
	req := pb.BusinessGetRequest{
		Id: g.Param("id"),
	}

	res, err := s.S.GetById(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *BusinessController) HandleGet(g *gin.Context) {
	q, err := strconv.Atoi(g.DefaultQuery("query", "0"))
	if err != nil {
		q = 0
	}
	req := pb.BusinessesGetRequest{
		CategoryId: g.Query("categoryId"),
		Zipcode:    g.Query("zipcode"),
		Offset:     g.DefaultQuery("offset", "0"),
		Limit:      g.DefaultQuery("limit", "0"),
		Mail:       g.Query("mail"),
		Phone:      g.Query("phone"),
		Query:      c.SORT_QUERY(q),
	}

	res, err := s.S.List(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *BusinessController) HandleNearGet(g *gin.Context) {
	req := pb.BusinessNearGetRequest{
		XUserId: g.GetString("userId"),
	}

	res, err := s.S.GetNear(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *BusinessController) HandleInterestGet(g *gin.Context) {
	req := pb.BusinessInterestGetRequest{}

	res, err := s.S.GetInterest(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *BusinessController) HandlePut(g *gin.Context) {
	req := pb.BusinessPutRequest{}
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

func (s *BusinessController) HandlePaymentMethodGet(g *gin.Context) {
	req := pb.BusinessPaymentMethodGetRequest{}
	req.XUserId = g.GetString("userId")
	res, err := s.S.GetPaymentMethod(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *BusinessController) HandlePaymentMethodPost(g *gin.Context) {
	req := pb.BusinessPaymentMethodPostRequest{}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.XUserId = g.GetString("userId")
	res, err := s.S.PostPaymentMethod(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *BusinessController) HandlePaymentMethodSetupPost(g *gin.Context) {
	req := pb.BusinessPaymentMethodSetupPostRequest{}
	req.XUserId = g.GetString("userId")
	res, err := s.S.SetupPaymentMethod(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *BusinessController) HandlePaymentMethodDeletePost(g *gin.Context) {
	req := pb.BusinessPaymentMethodDeletePostRequest{}
	req.XUserId = g.GetString("userId")
	res, err := s.S.DeletePaymentMethod(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *BusinessController) HandleServicesGet(g *gin.Context) {
	req := pb.BusinessServiceGetRequest{
		Id: g.Param("id"),
	}
	res, err := s.S.GetServices(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *BusinessController) HandleServicesPut(g *gin.Context) {
	req := pb.BusinessServicesPutRequest{}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.Id = g.Param("id")
	req.XUserId = g.GetString("userId")
	res, err := s.S.UpdateServices(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *BusinessController) TransactionsGet(g *gin.Context) {
	p, err := strconv.Atoi(g.DefaultQuery("query", "0"))
	if err != nil {
		p = 0
	}
	req := pb.BusinessTransactionsGetRequest{
		Query:  c.SORT_TRANSACTION(p),
		Limit:  g.DefaultQuery("limit", "5"),
		Offset: g.DefaultQuery("offset", "0"),
	}

	req.XUserId = g.GetString("userId")

	res, err := s.S.ListTransactions(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *BusinessController) TransactionsExport(g *gin.Context) {
	p, err := strconv.Atoi(g.DefaultQuery("query", "0"))
	if err != nil {
		p = 0
	}
	req := pb.BusinessTransactionsGetRequest{
		Query:  c.SORT_TRANSACTION(p),
		Limit:  g.DefaultQuery("limit", "5"),
		Offset: g.DefaultQuery("offset", "0"),
	}

	req.XUserId = g.GetString("userId")

	res, len, err := s.S.TransactionsExport(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	header := make(map[string]string)
	downloadName := time.Now().UTC().Format("data-20060102150405.xlsx")
	header["content-disposition"] = "attachment; filename=" + downloadName
	header["content-description"] = "File Transfer"
	g.DataFromReader(http.StatusOK, len, "application/octet-stream", res, header)
}

func (s BusinessController) HandleAdvertiseGet(g *gin.Context) {
	req := pb.AdvertiseGetRequest{
		Limit:  g.DefaultQuery("limit", "5"),
		Offset: g.DefaultQuery("offset", "0"),
	}

	req.XUserId = g.GetString("userId")

	res, err := s.S.GetAdvertisePackages(lib.ParseGinContext(g), &req)

	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)

}

func (s BusinessController) HandleAdvertiseDetailGet(g *gin.Context) {

	req := pb.AdvertiseDetailGetRequest{
		Id:     g.Param("id"),
		Limit:  g.DefaultQuery("limit", "5"),
		Offset: g.DefaultQuery("offset", "0"),
	}

	req.XUserId = g.GetString("userId")

	res, err := s.S.GetAdvertiseDetail(lib.ParseGinContext(g), &req)

	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)

}

func (s BusinessController) HandleAdvertiseOrdersGet(g *gin.Context) {
	req := pb.BusinessAdvertiseOrderGetRequest{
		Limit:  g.DefaultQuery("limit", "5"),
		Offset: g.DefaultQuery("offset", "0"),
	}

	req.XUserId = g.GetString("userId")

	res, err := s.S.GetAdvertiseOrder(lib.ParseGinContext(g), &req)

	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s BusinessController) HandleInvitationCodeGet(g *gin.Context) {
	req := pb.BusinessInvitationCodeGetRequest{}

	req.XUserId = g.GetString("userId")

	res, err := s.S.GetInvitationCode(lib.ParseGinContext(g), &req)

	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s BusinessController) HandleFreeContactGet(g *gin.Context) {
	req := pb.BusinessFreeContactGetRequest{
		Id: g.Param("id"),
	}

	res, err := s.S.GetNumberFreeContact(lib.ParseGinContext(g), &req)

	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *BusinessController) HandleValidateBuyAdvertisePost(g *gin.Context) {
	req := pb.BusinessValidateBuyAdvertisePostRequest{}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.XUserId = g.GetString("userId")
	res, err := s.S.ValidateBuyAdvertise(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *BusinessController) HandleBuyAdvertiseSetupPost(g *gin.Context) {
	req := pb.BusinessBuyAdvertiseSetupPostRequest{}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.XUserId = g.GetString("userId")

	res, err := s.S.BuyAdvertiseSetup(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *BusinessController) HandleBuyAdvertisePost(g *gin.Context) {
	req := pb.BusinessBuyAdvertisePostRequest{}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.XUserId = g.GetString("userId")
	res, err := s.S.BuyAdvertise(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *BusinessController) HandleVerifyRefCodePut(g *gin.Context) {
	req := pb.BusinessVerifyRefCodePutRequest{}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.Id = g.Param("Id")
	res, err := s.S.VerifyRefCode(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}
