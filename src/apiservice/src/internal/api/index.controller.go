package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"

	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/e"
	"github.com/aqaurius6666/apiservice/src/pb"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/webhook"
)

var indexSet = wire.NewSet(wire.Struct(new(IndexController), "*"), wire.Struct(new(IndexService), "*"))

type STRIPE_SIGNATURE_KEY string

type IndexController struct {
	S       IndexService
	SignKey STRIPE_SIGNATURE_KEY
}

func (s *IndexController) HandleIndexGet(g *gin.Context) {
	lib.Success(g, "Go go bruh bruh ...")
}

func (s *IndexController) HandleRandomGet(g *gin.Context) {
	rand := rand.Int()
	g.JSON(200, gin.H{"random": rand})
}

func (s *IndexController) HandleStripeKeyGet(g *gin.Context) {
	req := pb.StripeKeyGetRequest{}
	res, err := s.S.GetStripeKey(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *IndexController) HandlePaymentMethodGet(g *gin.Context) {
	req := pb.StripePaymentMethodGetRequest{}
	req.XUserId = g.GetString("userId")
	res, err := s.S.GetPaymentMethodInfo(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

// tao thanh toan
func (s *IndexController) HandleStripeSetupPost(g *gin.Context) {
	panic("not implemented")
	// req := pb.StripeSetupPostRequest{}
	// if err := lib.GetBody(g, &req); err != nil {
	// 	lib.BadRequest(g, err)
	// 	return
	// }
	// // res, err := s.S.SetupPayment(lib.ParseGinContext(g), &req)
	// if err != nil {
	// 	lib.BadRequest(g, err)
	// 	return
	// }
	// lib.Success(g, res)
}

func (s *IndexController) HandleSubscribePost(g *gin.Context) {
	req := pb.UnsubscribePostRequest{}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.XUserId = g.GetString("userId")
	res, err := s.S.UnsubscribeNotification(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *IndexController) HandleUnsubscribePost(g *gin.Context) {
	req := pb.SubscribePostRequest{}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.XUserId = g.GetString("userId")
	res, err := s.S.SubscribeNotification(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)

}

func (s *IndexController) HandleUploadUrlPost(g *gin.Context) {
	req := pb.UploadUrlPostRequest{}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.XUserId = g.GetString("userId")
	res, err := s.S.GetUploadUrl(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *IndexController) HandleWebHookPost(g *gin.Context) {
	var payload []byte
	body, ok := g.Get("rawBody")
	if !ok {
		return
	} else {
		payload = body.([]byte)
	}

	var stripeHeader string
	if len(g.Request.Header["Stripe-Signature"]) > 0 {
		stripeHeader = g.Request.Header["Stripe-Signature"][0]
	} else {
		fmt.Fprintf(os.Stderr, "%v", e.ErrStripeHeader)
		lib.BadRequest(g, e.ErrStripeHeader)
		return
	}
	event, err := webhook.ConstructEvent(payload, stripeHeader, string(s.SignKey))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		lib.BadRequest(g, err)
		return
	}
	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			return
		}
		err = s.S.UpdateTransactionWebHook(lib.ParseGinContext(g), paymentIntent)
		if err != nil {
			lib.BadRequest(g, err)
			return
		}
		lib.Success(g, nil)

	default:
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
	}
}

func (s *IndexController) HandleStatisticGet(g *gin.Context) {
	req := pb.StatisticGetRequest{}
	res, err := s.S.GetHomepageStatistic(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *IndexController) HandleCheckValidMailGet(g *gin.Context) {
	req := pb.ValidateMailGetRequest{
		Mail: g.Query("mail"),
	}
	res, err := s.S.ValidateMail(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}
