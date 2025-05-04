package controller

import (
	"TTCS/src/common/log"
	"TTCS/src/core/domain"
	"TTCS/src/core/service"
	"TTCS/src/present/httpui/request"
	"TTCS/src/present/httpui/response"
	"errors"
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

func (a *AuthController) SignUpOTP(ctx *gin.Context) {
	ctxReq := ctx.Request.Context()
	caller := "AuthController.GetOTP"

	email := ctx.Param("email")

	_, err := a.authService.SignUpOTP(ctxReq, email)
	if err != nil {
		log.Error(ctxReq, "[%v] failed to SignUpOTP: %+v", caller, err)
		a.ServeErrResponse(ctx, err)
		return
	}
	a.ServeSuccessResponse(ctx, true)
	return
}

func (a *AuthController) ResetOTP(ctx *gin.Context) {
	ctxReq := ctx.Request.Context()
	caller := "AuthController.GetOTP"

	email := ctx.Param("email")

	_, err := a.authService.ResetOTP(ctxReq, email)
	if err != nil {
		log.Error(ctxReq, "[%v] failed to SignUpOTP: %+v", caller, err)
		a.ServeErrResponse(ctx, err)
		return
	}
	a.ServeSuccessResponse(ctx, true)
	return
}

func (a *AuthController) ResetPassword(ctx *gin.Context) {
	ctxReq := ctx.Request.Context()
	caller := "AuthController.ResetPassword"

	var req request.ResetPasswordRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.Error(ctxReq, "[%v] failed to Bind: %+v", caller, err)
		a.ServeErrResponse(ctx, err)
		return
	}

	err := a.authService.ResetPassword(ctx, &req)
	if err != nil {
		log.Error(ctxReq, "[%v] failed to ResetPassword: %+v", caller, err)
		a.ServeErrResponse(ctx, err)
		return
	}
	a.ServeSuccessResponse(ctx, true)
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

	jwtToken, _, err := a.authService.SignUp(ctxReq, req)

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

	jwtToken, _, err := a.authService.Login(ctxReq, req)

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
	})
	return
}

func (a *AuthController) ChangePassword(ctx *gin.Context) {
	ctxReq := ctx.Request.Context()
	caller := "AuthController.ChangePassword"

	user, ok := ctx.Get("user")
	if !ok {
		err := errors.New("user not found")
		log.Error(ctxReq, "[%v] get user info error %+v", caller, err)
		a.ServeErrResponse(ctx, err)
		return
	}

	var req request.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error(ctxReq, "[%v] invalid param %+v", caller, err)
		a.ServeErrResponse(ctx, err)
		return
	}

	if err := a.authService.ChangePassword(ctxReq, user.(*domain.User).ID, req); err != nil {
		log.Error(ctxReq, "[%v] failed to change password: %+v", caller, err)
		a.ServeErrResponse(ctx, err)
		return
	}
	a.ServeSuccessResponse(ctx, true)
	return

}
