package controller

import (
	"TTCS/src/common/log"
	"TTCS/src/core/service"
	request2 "TTCS/src/present/httpui/request"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	*BaseController
	userService service.UserService
}

func NewUserController(baseController *BaseController, userService *service.UserService) *UserController {
	return &UserController{
		BaseController: baseController,
		userService:    *userService,
	}
}

func (u *UserController) UpdateInfo(ctx *gin.Context) {
	ctxReq := ctx.Request.Context()
	caller := "UserController.UpdateInfo"

	var req request2.UserInfo

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error(ctxReq, "[%v] invalid param %+v", caller, err)
		u.ServeErrResponse(ctx, err)
		return
	}

	id := ctx.Param("id")

	user, err := u.userService.Update(ctx, id, &req)
	if err != nil {
		log.Error(ctxReq, "[%v] update user info error %+v", caller, err)
		u.ServeErrResponse(ctx, err)
	}
	u.ServeSuccessResponse(ctx, user)
	return
}
