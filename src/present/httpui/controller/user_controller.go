package controller

import (
	"TTCS/src/common/log"
	"TTCS/src/core/domain"
	"TTCS/src/core/service"
	"TTCS/src/present/httpui/request"
	"TTCS/src/present/httpui/response"
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

	var req request.UserInfo

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
		return
	}
	u.ServeSuccessResponse(ctx, user)
	return
}

func (u *UserController) Create(ctx *gin.Context) {
	ctxReq := ctx.Request.Context()
	caller := "UserController.Create"

	var req request.UserInfo
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error(ctxReq, "[%v] invalid param %+v", caller, err)
		u.ServeErrResponse(ctx, err)
		return
	}
	user, err := u.userService.Create(ctx, &req)
	if err != nil {
		log.Error(ctxReq, "[%v] create user info error %+v", caller, err)
		u.ServeErrResponse(ctx, err)
		return
	}
	u.ServeSuccessResponse(ctx, response.UserFromDomain(user))
}

func (u *UserController) GetList(ctx *gin.Context) {
	ctxReq := ctx.Request.Context()
	caller := "UserController.GetList"

	var page request.Page
	if err := ctx.ShouldBindJSON(&page); err != nil {
		log.Error(ctxReq, "[%v] invalid param %+v", caller, err)
		u.ServeErrResponse(ctx, err)
		return
	}

	page.SetDefaults()

	users, err := u.userService.GetList(ctx, page)
	if err != nil {
		log.Error(ctxReq, "[%v] get user list error %+v", caller, err)
		u.ServeErrResponse(ctx, err)
		return
	}
	u.ServeSuccessResponse(ctx, response.UsersFromDomain(users))
}

func (u *UserController) GetDetail(ctx *gin.Context) {
	ctxReq := ctx.Request.Context()
	caller := "UserController.GetDetail"

	id := ctx.Param("id")

	user, err := u.userService.GetById(ctx, id)
	if err != nil {
		log.Error(ctxReq, "[%v] get user info error %+v", caller, err)
		u.ServeErrResponse(ctx, err)
		return
	}
	u.ServeSuccessResponse(ctx, response.UserFromDomain(user))
}

func (u *UserController) GetPayments(ctx *gin.Context) {
	ctxReq := ctx.Request.Context()
	caller := "UserController.GetPayment"

	user, ok := ctx.Get("user")
	if !ok {
		log.Error(ctxReq, "failed get user from context")
	}

	payments, err := u.userService.GetPayments(ctx, user.(domain.User).ID)

	if err != nil {
		log.Error(ctxReq, "[%v] get user payments error %+v", caller, err)
		u.ServeErrResponse(ctx, err)
		return
	}
	u.ServeSuccessResponse(ctx, payments)
}

func (u *UserController) GetOrders(ctx *gin.Context) {
	ctxReq := ctx.Request.Context()
	caller := "UserController.GetPayment"

	user, ok := ctx.Get("user")
	if !ok {
		log.Error(ctxReq, "failed get user from context")
	}

	orders, err := u.userService.GetOrders(ctx, user.(domain.User).ID)

	if err != nil {
		log.Error(ctxReq, "[%v] get user payments error %+v", caller, err)
		u.ServeErrResponse(ctx, err)
		return
	}
	u.ServeSuccessResponse(ctx, orders)
}
