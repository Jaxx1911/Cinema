package controller

import (
	"TTCS/src/common/log"
	"TTCS/src/common/ws"
	"TTCS/src/core/domain"
	"TTCS/src/core/service"
	"TTCS/src/present/httpui/request"
	"TTCS/src/present/httpui/response"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PaymentController struct {
	*BaseController
	paymentService *service.PaymentService
	hub            *ws.Hub
}

func NewPaymentController(baseController *BaseController, paymentService *service.PaymentService, hub *ws.Hub) *PaymentController {
	return &PaymentController{
		BaseController: baseController,
		paymentService: paymentService,
		hub:            hub,
	}
}

func (p *PaymentController) CallBack(ctx *gin.Context) {
	caller := "PaymentController.CallBack"
	ctxReq := ctx.Request.Context()

	var req request.PaymentCallback

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error(ctxReq, "[%v] invalid param %+v", caller, err)
		p.ServeErrResponse(ctx, err)
		return
	}

	payment, err := p.paymentService.HandleCallback(ctxReq, req)
	if err != nil {
		log.Error(ctxReq, "[%v] call back fail %+v", caller, err)
		p.ServeErrResponse(ctx, err)
		return
	}

	_ = p.hub.SendMessageToClient(*payment.UserID, ws.Message{
		Type: "payment",
		Data: Response{
			Key:     "payment",
			Body:    payment,
			Message: "success",
		},
	})
	log.Info(ctxReq, "[%v] payment success %+v", caller, payment)
	return
}

func (p *PaymentController) GetListByUserId(ctx *gin.Context) {
	caller := "PaymentController.GetListByUserId"
	ctxReq := ctx.Request.Context()

	user, ok := ctx.Get("user")
	if !ok {
		err := fmt.Errorf("failed to parse user in token")
		log.Error(ctxReq, "[%v] failed to get user from context %+v", caller, err)
		p.ServeErrResponse(ctx, err)
		return
	}

	payments, err := p.paymentService.GetPaymentsByUserID(ctxReq, user.(domain.User).ID)
	if err != nil {
		log.Error(ctxReq, "failed to get payments by user id %+v", user, err)
		p.ServeErrResponse(ctx, err)
		return
	}
	p.ServeSuccessResponse(ctx, response.ToPaymentsResponse(payments))

}

func (p *PaymentController) GetPaymentsByCinemaId(ctx *gin.Context) {
	caller := "PaymentController.GetPaymentsByCinemaId"
	ctxReq := ctx.Request.Context()

	cinemaId := ctx.Param("id")

	var req request.GetPaymentsByCinemaRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		log.Error(ctxReq, "[%v] invalid query parameters %+v", caller, err)
		p.ServeErrResponse(ctx, err)
		return
	}

	paymentDetails, err := p.paymentService.GetPaymentsByCinemaId(ctxReq, uuid.MustParse(cinemaId), req)
	if err != nil {
		log.Error(ctxReq, "[%v] failed to get payments by cinema id %+v", caller, err)
		p.ServeErrResponse(ctx, err)
		return
	}
	p.ServeSuccessResponse(ctx, response.ToPaymentCinemaDetailsResponse(paymentDetails))
}

func (p *PaymentController) AcceptAll(ctx *gin.Context) {
	caller := "PaymentController.AcceptAll"
	ctxReq := ctx.Request.Context()

	err := p.paymentService.AcceptAll(ctxReq)
	if err != nil {
		log.Error(ctxReq, "[%v] accept all payments failed, %+v", caller, err)
		p.ServeErrResponse(ctx, err)
		return
	}
	p.ServeSuccessResponse(ctx, nil)
}

func (p *PaymentController) GetList(ctx *gin.Context) {
	caller := "PaymentController.GetList"
	ctxReq := ctx.Request.Context()

	var req request.GetListPaymentRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		log.Error(ctxReq, "[%v] invalid query parameters %+v", caller, err)
		p.ServeErrResponse(ctx, err)
		return
	}

	payments, total, err := p.paymentService.GetList(ctxReq, req)
	if err != nil {
		log.Error(ctxReq, "[%v] failed to get payments list %+v", caller, err)
		p.ServeErrResponse(ctx, err)
		return
	}

	p.ServeSuccessResponse(ctx, response.MetaData{
		Data:       payments,
		TotalCount: total,
	})
}
