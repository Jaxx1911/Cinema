package controller

import (
	"TTCS/src/common/log"
	"TTCS/src/core/domain"
	"TTCS/src/core/service"
	"TTCS/src/present/httpui/request"
	"errors"
	"github.com/gin-gonic/gin"
)

type OrderController struct {
	*BaseController
	orderService *service.OrderService
}

func NewOrderController(baseController *BaseController, orderService *service.OrderService) *OrderController {
	return &OrderController{
		BaseController: baseController,
		orderService:   orderService,
	}
}

func (o OrderController) CreateOrder(ctx *gin.Context) {
	caller := "OrderController.CreateOrder"
	ctxReq := ctx.Request.Context()

	var req request.CreateOrderRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.Error(ctxReq, "[%v] invalid param %+v", caller, err)
		o.ServeErrResponse(ctx, err)
		return
	}
	user, ok := ctx.Get("user")
	if !ok {
		err := errors.New("user not found in context")
		log.Error(ctxReq, "[%v] not found user %+v", caller, err)
		o.ServeErrResponse(ctx, err)
		return
	}
	order, err := o.orderService.Create(ctxReq, user.(*domain.User).ID, req)
	if err != nil {
		log.Error(ctxReq, "[%v] create order failed, %+v", caller, err)
		o.ServeErrResponse(ctx, err)
		return
	}
	o.ServeSuccessResponse(ctx, order)
}
