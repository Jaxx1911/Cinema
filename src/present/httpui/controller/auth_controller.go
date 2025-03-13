package controller

import (
	"TTCS/src/common/log"
	"TTCS/src/core/service"
	"TTCS/src/present/httpui/request"
	"TTCS/src/present/httpui/response"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	*BaseController
	authService *service.AuthService
}

func NewAuthController(baseController *BaseController, authService *service.AuthService) *AuthController {
	return &AuthController{
		BaseController: baseController,
		authService:    authService,
	}
}

func (a *AuthController) GetOTP(ctx *gin.Context) {
	ctxReq := ctx.Request.Context()
	caller := "AuthController.GetOTP"

	email := ctx.Param("email")

	otp, err := a.authService.GenOTP(ctxReq, email)
	if err != nil {
		log.Error(ctxReq, "[%v] failed to GenOTP: %+v", caller, err)
		a.ServeErrResponse(ctx, err)
		return
	}
	a.ServeSuccessResponse(ctx, response.Otp{Otp: otp})
	return
}

func (a *AuthController) SignUp(ctx *gin.Context) {
	ctxReq := ctx.Request.Context()
	caller := "AuthController.SignUp"

	var req request.SignUpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error(ctxReq, "[%v] invalid param %+v", caller, err)
		a.ServeErrResponse(ctx, err)
		return
	}

	jwtToken, user, err := a.authService.SignUp(ctxReq, req)

	if err != nil {
		log.Error(ctxReq, "[%v] failed to register %+v", caller, err)
		a.ServeErrResponse(ctx, err)
		return
	}
	a.ServeSuccessResponse(ctx, response.LoginResp{
		Token: &response.Token{
			AccessToken:      jwtToken.AccessToken.Token,
			AccessExpiredAt:  jwtToken.AccessToken.Expire,
			RefreshToken:     jwtToken.RefreshToken.Token,
			RefreshExpiredAt: jwtToken.RefreshToken.Expire,
		},
		User: response.UserFromDomain(user),
	})
	return
}

func (a *AuthController) Login(ctx *gin.Context) {
	ctxReq := ctx.Request.Context()
	caller := "AuthController.Login"

	var req request.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error(ctxReq, "[%v] invalid param %+v", caller, err)
		a.ServeErrResponse(ctx, err)
		return
	}

	jwtToken, user, err := a.authService.Login(ctxReq, req)

	if err != nil {
		log.Error(ctxReq, "[%v] failed to create token", caller)
		a.ServeErrResponse(ctx, err)
		return
	}
	a.ServeSuccessResponse(ctx, response.LoginResp{
		Token: &response.Token{
			AccessToken:      jwtToken.AccessToken.Token,
			AccessExpiredAt:  jwtToken.AccessToken.Expire,
			RefreshToken:     jwtToken.RefreshToken.Token,
			RefreshExpiredAt: jwtToken.RefreshToken.Expire,
		},
		User: response.UserFromDomain(user),
	})
	return
}

func (a *AuthController) ChangePassword(ctx *gin.Context) {
	ctxReq := ctx.Request.Context()
	caller := "AuthController.ChangePassword"

	var req request.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error(ctxReq, "[%v] invalid param %+v", caller, err)
		a.ServeErrResponse(ctx, err)
		return
	}

	if err := a.authService.ChangePassword(ctxReq, req); err != nil {
		log.Error(ctxReq, "[%v] failed to change password: %+v", caller, err)
		a.ServeErrResponse(ctx, err)
		return
	}
	a.ServeSuccessResponse(ctx, true)
	return

}
