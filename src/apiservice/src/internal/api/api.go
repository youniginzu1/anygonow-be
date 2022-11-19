package api

import (
	"github.com/aqaurius6666/apiservice/src/internal/api/middleware"
	"github.com/aqaurius6666/apiservice/src/internal/db"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var ApiServerSet = wire.NewSet(wire.Struct(new(ApiServer), "*"),
	gin.New,
	indexSet,
	authSet,
	businessSet,
	userSet,
	contactSet,
	middleware.Set,
	adminSet,
	categorySet,
	orderSet,
	feedbackSet,
	chatSet,
)

type ApiServer struct {
	G          *gin.Engine
	Logger     *logrus.Logger
	ServerRepo db.ServerRepo
	Index      IndexController
	Business   BusinessController
	Contact    ContactController
	Auth       AuthController
	Admin      AdminController
	Category   CategoryController
	Mid        middleware.Middleware
	Users      UserController
	Order      OrderController
	Feedback   FeedbackController
	Chat       ChatController
}

func (s *ApiServer) RegisterEndpoint() {
	// gin.SetMode(gin.ReleaseMode)
	s.G.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "PUT", "POST", "DELETE"},
		AllowHeaders:  []string{"Authorization", "Content-Type", "User-Agent"},
		ExposeHeaders: []string{"content-disposition", "content-description"},
	}))
	s.G.Use(otelgin.Middleware(c.SERVICE_NAME))
	s.G.Use(gin.Recovery())
	s.G.Use(s.Mid.HandleBody())
	s.G.Use(s.Mid.LoggerMid())

	s.G.GET("/", s.Index.HandleIndexGet)

	api := s.G.Group("/api")
	api.GET("", s.Index.HandleIndexGet)
	api.POST("/web-hook", s.Index.HandleWebHookPost)
	api.GET("/stripe/key", s.Index.HandleStripeKeyGet)
	api.POST("/stripe/setup", s.Mid.CheckAuth, s.Index.HandleStripeSetupPost)
	api.GET("/stripe/payment-method", s.Mid.CheckAuth, s.Index.HandlePaymentMethodGet)
	api.GET("/", s.Index.HandleIndexGet)
	api.GET("/random", s.Index.HandleRandomGet)
	api.GET("/statistics", s.Index.HandleStatisticGet)
	api.GET("/check-valid-mail", s.Index.HandleCheckValidMailGet)

	api.POST("/upload-url", s.Mid.CheckAuth, s.Index.HandleUploadUrlPost)
	api.POST("/subscribe", s.Mid.CheckAuth, s.Index.HandleSubscribePost)
	api.POST("/unsubscribe", s.Mid.CheckAuth, s.Index.HandleUnsubscribePost)

	authGroup := api.Group("/auth")
	authGroup.POST("/credential", s.Auth.HandleCredentialPost)
	authGroup.POST("/mail", s.Mid.CheckAuth, s.Auth.HandleMailPost)
	authGroup.POST("/password", s.Mid.CheckAuth, s.Auth.HandlePasswordPost)
	authGroup.GET("/check", s.Auth.HandleCheckGet)
	authGroup.POST("/otp", s.Auth.HandleOTPPost)
	authGroup.POST("/otp/resend", s.Auth.HandleResendOTPPost)
	authGroup.POST("/forgot", s.Auth.HandleForgotPost)
	authGroup.POST("/forgot/reset", s.Auth.HandleForgotResetPost)
	authGroup.POST("/ping", s.Mid.CheckAuth, s.Auth.HandlePingPost)
	authGroup.PUT("/change-mail-and-pass", s.Mid.CheckAuth, s.Auth.HandleMailPassPut)

	userGroup := api.Group("/users")
	userGroup.POST("", s.Users.HandleUserPost)
	userGroup.GET("/:id", s.Mid.CheckAuth, s.Users.HandleGetById)
	userGroup.PUT("/:id", s.Mid.CheckAuth, s.Users.HandlePut)
	userGroup.GET("/state", s.Users.HandleStateGet)

	businessGroup := api.Group("/businesses")
	businessGroup.GET("/payment-summary", s.Mid.CheckAuth, s.Mid.Only(c.ROLE_HANDYMAN), s.Business.TransactionsGet)
	businessGroup.GET("/payment-summary-export", s.Mid.CheckAuth, s.Mid.Only(c.ROLE_HANDYMAN), s.Business.TransactionsExport)
	businessGroup.POST("", s.Business.HandlePost)
	businessGroup.GET("", s.Business.HandleGet)
	businessGroup.GET("/:id/rating", s.Business.HandleRatingGet)
	businessGroup.GET("/:id/feedbacks", s.Business.HandleFeedbackGet)
	businessGroup.GET("/:id/services", s.Business.HandleServicesGet)
	businessGroup.GET("/:id/free-contact", s.Mid.CheckAuth, s.Mid.Only(c.ROLE_HANDYMAN), s.Business.HandleFreeContactGet)
	businessGroup.PUT("/:id/services", s.Mid.CheckAuth, s.Mid.Only(c.ROLE_HANDYMAN), s.Business.HandleServicesPut)
	businessGroup.PUT("/:id/verify-refcode", s.Mid.CheckAuth, s.Mid.Only(c.ROLE_HANDYMAN), s.Business.HandleVerifyRefCodePut)
	businessGroup.PUT("/:id", s.Mid.CheckAuth, s.Business.HandlePut)
	businessGroup.GET("/interest", s.Business.HandleInterestGet)
	businessGroup.GET("/near", s.Mid.CheckAuth, s.Business.HandleNearGet)
	businessGroup.GET("/invitation-code", s.Mid.CheckAuth, s.Mid.Only(c.ROLE_HANDYMAN), s.Business.HandleInvitationCodeGet)
	businessGroup.GET("/payment-method", s.Mid.CheckAuth, s.Business.HandlePaymentMethodGet)
	businessGroup.POST("/payment-method", s.Mid.CheckAuth, s.Business.HandlePaymentMethodPost)
	businessGroup.POST("/payment-method/setup", s.Mid.CheckAuth, s.Business.HandlePaymentMethodSetupPost)
	businessGroup.POST("/payment-method/delete", s.Mid.CheckAuth, s.Business.HandlePaymentMethodDeletePost)
	businessGroup.GET("/:id", s.Business.HandleGetById)
	businessGroup.GET("/promote", s.Mid.CheckAuth, s.Mid.Only(c.ROLE_HANDYMAN), s.Business.HandleAdvertiseGet)
	businessGroup.GET("/promote/:id", s.Mid.CheckAuth, s.Mid.Only(c.ROLE_HANDYMAN), s.Business.HandleAdvertiseDetailGet)
	businessGroup.GET("/promote/order", s.Mid.CheckAuth, s.Mid.Only(c.ROLE_HANDYMAN), s.Business.HandleAdvertiseOrdersGet)
	businessGroup.POST("/buy-promote/setup", s.Mid.CheckAuth, s.Mid.Only(c.ROLE_HANDYMAN), s.Business.HandleBuyAdvertiseSetupPost)
	businessGroup.POST("/buy-promote", s.Mid.CheckAuth, s.Mid.Only(c.ROLE_HANDYMAN), s.Business.HandleBuyAdvertisePost)
	businessGroup.POST("/buy-promote/validate", s.Mid.CheckAuth, s.Mid.Only(c.ROLE_HANDYMAN), s.Business.HandleValidateBuyAdvertisePost)

	contactGroup := api.Group("/contacts")
	contactGroup.GET("/states", s.Contact.HandleListStates)
	contactGroup.GET("/state/:id", s.Contact.HandleGetState)
	contactGroup.GET("/:id", s.Mid.CheckAuth, s.Contact.HandleGetById)
	contactGroup.PUT("/:id", s.Mid.CheckAuth, s.Contact.HandlePut)

	adminGroup := api.Group("/admin")
	adminGroup.GET("/categories", s.Mid.CheckAuth, s.Mid.OnlyAdmin(), s.Admin.HandleCategoriesGet)
	adminGroup.POST("/categories", s.Mid.CheckAuth, s.Mid.OnlyAdmin(), s.Admin.HandleCategoriesPost)
	adminGroup.POST("/categories/delete", s.Mid.CheckAuth, s.Mid.OnlyAdmin(), s.Admin.HandleCategoriesPostDelete)
	adminGroup.POST("/categories/edit", s.Mid.CheckAuth, s.Mid.OnlyAdmin(), s.Admin.HandleCategoriesPostEdit)
	adminGroup.GET("/groups", s.Mid.CheckAuth, s.Mid.OnlyAdmin(), s.Admin.HandleGroupsGet)
	adminGroup.POST("/groups", s.Mid.CheckAuth, s.Mid.OnlyAdmin(), s.Admin.HandleGroupsPost)
	adminGroup.PUT("/groups/:id", s.Mid.CheckAuth, s.Mid.OnlyAdmin(), s.Admin.HandleGroupsPut)
	adminGroup.POST("/users/:id/ban", s.Mid.CheckAuth, s.Mid.OnlyAdmin(), s.Admin.HandleUserBanPost)
	adminGroup.POST("/users/:id/unban", s.Mid.CheckAuth, s.Mid.OnlyAdmin(), s.Admin.HandleUserUnbanPost)
	adminGroup.POST("/users/:id/delete", s.Mid.CheckAuth, s.Mid.OnlyAdmin(), s.Admin.HandleUserDeletePost)
	adminGroup.GET("/businesses", s.Mid.CheckAuth, s.Mid.OnlyAdmin(), s.Admin.HandleBusinessesGet)
	adminGroup.GET("/users", s.Mid.CheckAuth, s.Mid.OnlyAdmin(), s.Admin.HandleUsersGet)
	adminGroup.POST("/businesses/:id/ban", s.Mid.CheckAuth, s.Mid.OnlyAdmin(), s.Admin.HandleBusinessBanPost)
	adminGroup.POST("/businesses/:id/unban", s.Mid.CheckAuth, s.Mid.OnlyAdmin(), s.Admin.HandleBusinessUnbanPost)
	adminGroup.POST("/businesses/:id/delete", s.Mid.CheckAuth, s.Mid.OnlyAdmin(), s.Admin.HandleBusinessDeletePost)
	adminGroup.POST("/promote-management", s.Mid.CheckAuth, s.Mid.OnlyAdmin(), s.Admin.HandleAdvertiseManagementPost)
	adminGroup.GET("/promote-management", s.Mid.CheckAuth, s.Mid.OnlyAdmin(), s.Admin.HandleAdvertiseManagementGet)
	adminGroup.PUT("/promote-management/:id", s.Mid.CheckAuth, s.Mid.OnlyAdmin(), s.Admin.HandleAdvertiseManagementPut)
	adminGroup.POST("/promote-management/:id/delete", s.Mid.CheckAuth, s.Mid.OnlyAdmin(), s.Admin.HandleAdvertiseManagementDeletePost)

	categoryGroup := api.Group("/categories")
	categoryGroup.GET("", s.Category.HandleGet)
	categoryGroup.GET("/:id", s.Category.HandleByIdGet)

	orderGroup := api.Group("/orders")
	orderGroup.GET("/projects", s.Mid.CheckAuth, s.Mid.Only(c.ROLE_CUSTOMER), s.Order.HandleProjectsGet)
	orderGroup.GET("/already-ordered", s.Mid.CheckAuth, s.Mid.Only(c.ROLE_CUSTOMER), s.Order.HandleBusinessesAlreadyOrderedGet)
	orderGroup.POST("/projects/cancel", s.Mid.CheckAuth, s.Mid.Only(c.ROLE_CUSTOMER), s.Order.HandleProjectCancelPost)
	orderGroup.POST("", s.Mid.CheckAuth, s.Mid.Only(c.ROLE_CUSTOMER), s.Order.HandlePost)
	orderGroup.POST("/connect", s.Mid.CheckAuth, s.Mid.Only(c.ROLE_HANDYMAN), s.Order.HandleConnectPost)
	orderGroup.POST("/connect-all", s.Mid.CheckAuth, s.Mid.Only(c.ROLE_HANDYMAN), s.Order.HandleConnectAllPost)
	orderGroup.POST("/reject", s.Mid.CheckAuth, s.Mid.Only(c.ROLE_HANDYMAN), s.Order.HandleRejectPost) // must be pending or connnected
	orderGroup.POST("/cancel", s.Mid.CheckAuth, s.Mid.Only(c.ROLE_CUSTOMER), s.Order.HandleCancelPost) // must be pending
	orderGroup.POST("/complete", s.Mid.CheckAuth, s.Mid.Only(c.ROLE_CUSTOMER, c.ROLE_HANDYMAN), s.Order.HandleCompletePost) // must be connected
	orderGroup.GET("", s.Mid.CheckAuth, s.Order.HandleGet)

	feedbackGroup := api.Group("/feedbacks")
	feedbackGroup.POST("", s.Mid.CheckAuth, s.Mid.Only(c.ROLE_CUSTOMER), s.Feedback.HandlePost)
	feedbackGroup.GET("", s.Mid.CheckAuth, s.Mid.Only(c.ROLE_CUSTOMER), s.Feedback.HandleGet)
	feedbackGroup.PUT("/:id", s.Mid.CheckAuth, s.Mid.Only(c.ROLE_CUSTOMER), s.Feedback.HandlePut)

	chatGroup := api.Group("/chatservice")
	chatGroup.POST("/conversations", s.Mid.CheckAuth, s.Chat.HandleConversationsGet)

}
