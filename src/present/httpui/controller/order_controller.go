package controller

import (
	"TTCS/src/common/genqr"
	"TTCS/src/common/log"
	"TTCS/src/core/domain"
	"TTCS/src/core/service"
	"TTCS/src/present/httpui/request"
	"TTCS/src/present/httpui/response"
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
	o.ServeSuccessResponse(ctx, response.ToOrderResponse(*order))
}

func (o OrderController) GetOrderDetailsWithQr(ctx *gin.Context) {
	caller := "OrderController.GetOrderDetails"
	ctxReq := ctx.Request.Context()

	id := ctx.Param("id")
	order, err := o.orderService.GetById(ctxReq, id)
	if err != nil {
		log.Error(ctxReq, "[%v] get order details failed, %+v", caller, err)
		o.ServeErrResponse(ctx, err)
		return
	}
	qrText := genqr.QrGenerator.GenerateQrCode(int(order.TotalPrice), order.ID.String())
	o.ServeSuccessResponse(ctx, response.OrderWithQrResponse{
		Order:  response.ToOrderResponse(*order),
		QrText: qrText,
	})
}
