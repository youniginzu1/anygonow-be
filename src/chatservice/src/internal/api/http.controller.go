package api

import (
	"github.com/aqaurius6666/chatservice/src/internal/lib"
	"github.com/aqaurius6666/chatservice/src/internal/var/e"
	"github.com/aqaurius6666/chatservice/src/pb/chatpb"
	"github.com/aqaurius6666/chatservice/src/services/twilloclient"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

var (
	httpSet = wire.NewSet(wire.Struct(new(HTTPController), "*"), wire.Struct(new(HTTPService), "*"))
)

type HTTPController struct {
	Logger *logrus.Logger
	S      *HTTPService
	Twillo *twilloclient.TwilloClient
}

func (s *HTTPController) HandleTwilloCallback(g *gin.Context) {
	message, ok := g.Get("twillo-message")
	if !ok {
		lib.BadRequest(g, e.ErrBodyInvalid)
		return
	}
	// raw, ok := g.Get("raw-body")
	// if !ok {
	// 	lib.BadRequest(g, e.ErrBodyInvalid)
	// 	return
	// }
	messageCallback, ok := message.(twilloclient.TwilloMessageCallbackData)
	if !ok {
		lib.BadRequest(g, e.ErrBodyInvalid)
		return
	}
	// signature := g.GetHeader("X-Twilio-Signature")
	// if !s.Twillo.VerifySignature(signature, raw.([]byte)) {
	// 	lib.BadRequest(g, e.ErrBodyInvalid)
	// 	return
	// }
	res, err := s.S.ProcessMessageCallback(g.Request.Context(), &messageCallback)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *HTTPController) HandleFetchPost(g *gin.Context) {
	req := chatpb.FetchPostRequest{}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	res, err := s.S.FetchChats(g.Request.Context(), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *HTTPController) HandleFetchByIdPost(g *gin.Context) {
	req := chatpb.FetchOnePostRequest{}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.Id = g.Param("id")
	req.XUserId = g.GetString("userId")

	res, err := s.S.FetchChatById(g.Request.Context(), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *HTTPController) HandleSendPost(g *gin.Context) {
	req := chatpb.SendPostRequest{}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.Id = g.Param("id")
	req.XUserId = g.GetString("userId")

	res, err := s.S.SendChat(g.Request.Context(), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *HTTPController) HandlePostConversation(g *gin.Context) {
	req := chatpb.ConversationPostRequest{}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.XUserId = g.GetString("userId")

	res, err := s.S.GetConversation(g.Request.Context(), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *HTTPController) HandleGetNotification(g *gin.Context) {
	req := chatpb.NotificationGetRequest{}

	req.XUserId = g.GetString("userId")

	res, err := s.S.GetNotification(g.Request.Context(), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s *HTTPController) HandlePutNotification(g *gin.Context) {
	req := chatpb.NotificationSeenPutRequest{}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.XUserId = g.GetString("userId")

	res, err := s.S.SeenNotification(g.Request.Context(), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}
