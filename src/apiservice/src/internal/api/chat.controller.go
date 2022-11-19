package api

import (
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/pb"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var chatSet = wire.NewSet(wire.Struct(new(ChatController), "*"), wire.Struct(new(ChatService), "*"))

type ChatController struct {
	S ChatService
}

func (s *ChatController) HandleConversationsGet(g *gin.Context) {

	req := pb.ConversationPostRequest{}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.XUserId = g.GetString("userId")
	res, err := s.S.GetConversations(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}
