package controller

import (
	"TTCS/src/common/log"
	"TTCS/src/core/domain"
	"TTCS/src/core/service"
	"TTCS/src/present/httpui/request"
	"TTCS/src/present/httpui/response"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
func (u *UserController) Update(ctx *gin.Context) {
	ctxReq := ctx.Request.Context()
	caller := "UserController.Update"

	id := ctx.Param("id")

	var req request.UserInfo

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error(ctxReq, "[%v] invalid param %+v", caller, err)
		u.ServeErrResponse(ctx, err)
		return
	}

	user, err := u.userService.Update(ctx, uuid.MustParse(id), &req)
	if err != nil {
		log.Error(ctxReq, "[%v] update user info error %+v", caller, err)
		u.ServeErrResponse(ctx, err)
		return
	}
	u.ServeSuccessResponse(ctx, response.UserFromDomain(user))
	return
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

	user, ok := ctx.Get("user")
	if !ok {
		err := errors.New("user not found")
		log.Error(ctxReq, "[%v] get user info error %+v", caller, err)
		u.ServeErrResponse(ctx, err)
		return
	}

	user, err := u.userService.Update(ctx, user.(*domain.User).ID, &req)
	if err != nil {
		log.Error(ctxReq, "[%v] update user info error %+v", caller, err)
		u.ServeErrResponse(ctx, err)
		return
	}
	u.ServeSuccessResponse(ctx, response.UserFromDomain(user.(*domain.User)))
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

	var req request.GetListUser
	if err := ctx.ShouldBindQuery(&req); err != nil {
		log.Error(ctxReq, "[%v] failed to bind request: %v", caller, err)
		u.ServeErrResponse(ctx, err)
		return
	}

	// Set default values for pagination
	req.SetDefaults()

	users, total, err := u.userService.GetList(ctxReq, req)
	if err != nil {
		log.Error(ctxReq, "[%v] failed to get user list: %v", caller, err)
		u.ServeErrResponse(ctx, err)
		return
	}

	u.ServeSuccessResponse(ctx, gin.H{
		"users": response.UsersFromDomain(users),
		"total": total,
		"page":  req.Page,
		"limit": req.Limit,
	})
}

func (u *UserController) GetDetail(ctx *gin.Context) {
	ctxReq := ctx.Request.Context()
	caller := "UserController.GetDetail"

	user, ok := ctx.Get("user")
	if !ok {
		err := errors.New("user not found")
		log.Error(ctxReq, "[%v] get user info error %+v", caller, err)
		u.ServeErrResponse(ctx, err)
		return
	}

	u.ServeSuccessResponse(ctx, response.UserFromDomain(user.(*domain.User)))
}

func (u *UserController) GetPayments(ctx *gin.Context) {
	ctxReq := ctx.Request.Context()
	caller := "UserController.GetPayment"

	user, ok := ctx.Get("user")
	if !ok {
		err := errors.New("user not found in context")
		log.Error(ctxReq, "[%v] failed get user from context %+v", caller, err)
		u.ServeErrResponse(ctx, err)
	}

	payments, err := u.userService.GetPayments(ctx, user.(*domain.User).ID)

	if err != nil {
		log.Error(ctxReq, "[%v] get user payments error %+v", caller, err)
		u.ServeErrResponse(ctx, err)
		return
	}
	u.ServeSuccessResponse(ctx, response.ToPaymentsResponse(payments))
}

func (u *UserController) GetOrders(ctx *gin.Context) {
	ctxReq := ctx.Request.Context()
	caller := "UserController.GetPayment"

	user, ok := ctx.Get("user")
	if !ok {
		err := errors.New("user not found in context")
		log.Error(ctxReq, "[%v] failed get user from context %+v", caller, err)
		u.ServeErrResponse(ctx, err)
	}

	orders, err := u.userService.GetOrders(ctx, user.(*domain.User).ID)

	if err != nil {
		log.Error(ctxReq, "[%v] get user payments error %+v", caller, err)
		u.ServeErrResponse(ctx, err)
		return
	}
	u.ServeSuccessResponse(ctx, response.ToOrdersResponse(orders))
}

func (u *UserController) ChangeAvatar(ctx *gin.Context) {
	ctxReq := ctx.Request.Context()
	caller := "UserController.ChangeAvatar"

	user, ok := ctx.Get("user")
	if !ok {
		err := errors.New("user not found in context")
		log.Error(ctxReq, "[%v] failed get user from context %+v", caller, err)
		u.ServeErrResponse(ctx, err)
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		log.Error(ctxReq, "[%v] failed get file from form param %+v", caller, err)
		u.ServeErrResponse(ctx, err)
		return
	}

	nUser, err := u.userService.ChangeAvatar(ctx, file, user.(*domain.User))
	if err != nil {
		log.Error(ctxReq, "[%v] change user avatar error %+v", caller, err)
		u.ServeErrResponse(ctx, err)
		return
	}
	u.ServeSuccessResponse(ctx, nUser)
}

func (u *UserController) GetMe(ctx *gin.Context) {
	ctxReq := ctx.Request.Context()
	caller := "UserController.GetMe"

	user, ok := ctx.Get("user")
	if !ok {
		err := errors.New("user not found in context")
		log.Error(ctxReq, "[%v] failed to get user from context: %v", caller, err)
		u.ServeErrResponse(ctx, err)
		return
	}

	u.ServeSuccessResponse(ctx, response.UserFromDomain(user.(*domain.User)))
}

func (u *UserController) Delete(ctx *gin.Context) {
	ctxReq := ctx.Request.Context()
	caller := "UserController.Delete"

	id := ctx.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		log.Error(ctxReq, "[%v] invalid user id: %v", caller, err)
		u.ServeErrResponse(ctx, err)
		return
	}

	if err := u.userService.Delete(ctxReq, userID); err != nil {
		log.Error(ctxReq, "[%v] failed to delete user: %v", caller, err)
		u.ServeErrResponse(ctx, err)
		return
	}

	u.ServeSuccessResponse(ctx, true)
}
