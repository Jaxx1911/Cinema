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

func (c *AuthController) Login(ctx *gin.Context) {
	ctxReq := ctx.Request.Context()
	caller := "AuthController.Login"

	var req request.LoginRequest
	if err := c.BindAndValidateRequest(ctx, req); err != nil {
		log.Error(ctxReq, "[%v] invalid param %+v", caller, err)
		c.Error(ctx, err)
		return
	}

	jwtToken, user, err := c.authService.Login(ctxReq, req)

	if err != nil {
		log.Error(ctxReq, "[%v] failed to create token", caller)
		c.Error(ctx, err)
		return
	}
	c.Success(ctx, response.LoginResp{
		Token: jwtToken,
		User:  response.UserFromDomain(user),
	})

}
