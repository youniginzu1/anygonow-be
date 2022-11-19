package api

import (
	"github.com/aqaurius6666/chatservice/src/internal/api/middleware"
	"github.com/aqaurius6666/chatservice/src/internal/db"
	"github.com/aqaurius6666/chatservice/src/internal/model"
	"github.com/aqaurius6666/chatservice/src/internal/var/c"
	"github.com/aqaurius6666/chatservice/src/pb/chatpb"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var ApiServerSet = wire.NewSet(wire.Struct(new(ApiServer), "*"), httpSet, middleware.Set, indexSet)

type ApiServer struct {
	chatpb.UnimplementedChatServiceServer `wire:"-"`
	G                                     *gin.Engine
	Logger                                *logrus.Logger
	Model                                 model.Server
	Repo                                  db.ServerRepo
	HTTP                                  *HTTPController
	Mid                                   middleware.Middleware
	Index                                 *IndexController
	HTTPService                           *HTTPService
}

func (s *ApiServer) RegisterEndpoint() {
	gin.SetMode(gin.DebugMode)
	s.G.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "PUT", "POST", "DELETE"},
		AllowHeaders: []string{"Authorization", "Content-Type", "User-Agent"},
	}))
	s.G.Use(otelgin.Middleware(c.SERVICE_NAME))
	s.G.Use(gin.Recovery())
	s.G.Use(s.Mid.HandleBody())
	s.G.Use(s.Mid.LoggerMid())

	s.G.GET("/", s.Index.HandleIndexGet)
	api := s.G.Group("/api")
	chat := api.Group("/chat")
	chat.GET("", s.Index.HandleIndexGet)
	chat.GET("/", s.Index.HandleIndexGet)
	chat.POST("/twillo/callback", s.HTTP.HandleTwilloCallback)

	chat.POST("/conversation", s.Mid.CheckAuth, s.HTTP.HandlePostConversation)
	chat.GET("/notification", s.Mid.CheckAuth, s.HTTP.HandleGetNotification)
	chat.PUT("/notification", s.Mid.CheckAuth, s.HTTP.HandlePutNotification)
	chat.POST("/:id/fetch", s.Mid.CheckAuth, s.HTTP.HandleFetchByIdPost)
	chat.POST("/:id/send", s.Mid.CheckAuth, s.HTTP.HandleSendPost)

}
