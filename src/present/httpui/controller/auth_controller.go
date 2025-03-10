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
	userService *service.UserService
}

func NewAuthController(baseController *BaseController, authService *service.AuthService, userService *service.UserService) *AuthController {
	return &AuthController{
		BaseController: baseController,
		authService:    authService,
		userService:    userService,
	}
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
		Token: jwtToken,
		User:  response.UserFromDomain(user),
	})

}
