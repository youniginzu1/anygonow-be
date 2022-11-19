package api

import (
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/pb"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var authSet = wire.NewSet(wire.Struct(new(AuthController), "*"), wire.Struct(new(AuthService), "*"))

type AuthController struct {
	S AuthService
}

func (s AuthController) HandleCheckGet(g *gin.Context) {
	req := pb.AuthCheckGetRequest{
		Identifier: g.DefaultQuery("identifier", ""),
	}
	res, err := s.S.CheckCredential(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s AuthController) HandleCredentialPost(g *gin.Context) {
	req := pb.AuthCredentialRequest{}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	res, err := s.S.GetCredential(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s AuthController) HandleOTPPost(g *gin.Context) {
	req := pb.AuthOTPPostRequest{}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	res, err := s.S.VerifyOTP(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s AuthController) HandleResendOTPPost(g *gin.Context) {
	req := pb.AuthResendOTPPostRequest{}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	res, err := s.S.ResendOTP(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s AuthController) HandleMailPost(g *gin.Context) {
	req := pb.AuthMailPostRequest{}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.XUserId = g.GetString("userId")
	res, err := s.S.ChangeMail(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s AuthController) HandlePingPost(g *gin.Context) {
	req := pb.AuthPingRequest{
		XUserId: g.GetString("userId"),
		XRole:   lib.MustGetRole(g),
	}
	res, err := s.S.Ping(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s AuthController) HandlePasswordPost(g *gin.Context) {
	req := pb.AuthPasswordPostRequest{}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.XUserId = g.GetString("userId")

	res, err := s.S.ChangePassword(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s AuthController) HandleForgotPost(g *gin.Context) {
	req := pb.AuthForgotPostRequest{}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	res, err := s.S.ForgotPassword(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s AuthController) HandleForgotResetPost(g *gin.Context) {
	req := pb.AuthForgotResetPostRequest{}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	res, err := s.S.ResetPassword(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}
	lib.Success(g, res)
}

func (s AuthController) HandleMailPassPut(g *gin.Context) {
	req := pb.AuthChangeMailAndPassPostRequest{}
	if err := lib.GetBody(g, &req); err != nil {
		lib.BadRequest(g, err)
		return
	}
	req.XUserId = g.GetString("userId")
	res, err := s.S.ChangeMailAndPass(lib.ParseGinContext(g), &req)
	if err != nil {
		lib.BadRequest(g, err)
		return
	}

	lib.Success(g, res)
}
