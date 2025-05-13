package controller

import (
	"TTCS/src/common/log"
	"TTCS/src/core/service"
	"TTCS/src/present/httpui/request"
	"TTCS/src/present/httpui/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DiscountController struct {
	*BaseController
	discountService *service.DiscountService
}

func NewDiscountController(baseController *BaseController, discountService *service.DiscountService) *DiscountController {
	return &DiscountController{
		BaseController:  baseController,
		discountService: discountService,
	}
}

func (d DiscountController) GetDiscountByCode(ctx *gin.Context) {
	caller := "DiscountController.GetDiscountByCode"
	ctxReq := ctx.Request.Context()

	code := ctx.Param("code")
	discount, err := d.discountService.GetDiscountByCode(ctxReq, code)
	if err != nil {
		log.Error(ctxReq, "[%v] failed to get discount %+v", caller, err)
		d.ServeErrResponse(ctx, err)
		return
	}
	d.ServeSuccessResponse(ctx, response.ToDiscountResponse(*discount))
	return
}

func (d DiscountController) GetDiscounts(ctx *gin.Context) {
	caller := "DiscountController.GetDiscounts"
	ctxReq := ctx.Request.Context()

	discounts, err := d.discountService.GetListDiscount(ctxReq)
	if err != nil {
		log.Error(ctxReq, "[%v] failed to get discounts %+v", caller, err)
		d.ServeErrResponse(ctx, err)
		return
	}
	d.ServeSuccessResponse(ctx, response.ToListDiscountResponse(discounts))
}

func (d DiscountController) GetDiscountByID(ctx *gin.Context) {
	caller := "DiscountController.GetDiscountByID"
	ctxReq := ctx.Request.Context()
	id := ctx.Param("id")

	discount, err := d.discountService.GetDiscount(ctxReq, uuid.MustParse(id))
	if err != nil {
		log.Error(ctxReq, "[%v] failed to get discount %+v", caller, err)
		d.ServeErrResponse(ctx, err)
		return
	}
	d.ServeSuccessResponse(ctx, response.ToDiscountResponse(*discount))
}

func (d DiscountController) Create(ctx *gin.Context) {
	caller := "DiscountController.Create"
	ctxReq := ctx.Request.Context()

	var req request.Discount
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error(ctxReq, "[%v] failed to parse request body %+v", caller, err)
		d.ServeErrResponse(ctx, err)
		return
	}
	discount, err := d.discountService.Create(ctxReq, req)
	if err != nil {
		log.Error(ctxReq, "[%v] failed to create discount %+v", caller, err)
		d.ServeErrResponse(ctx, err)
		return
	}
	d.ServeSuccessResponse(ctx, response.ToDiscountResponse(*discount))
}

func (d DiscountController) Update(ctx *gin.Context) {
	caller := "DiscountController.Update"
	ctxReq := ctx.Request.Context()

	id := ctx.Param("id")
	var req request.Discount
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error(ctxReq, "[%v] failed to parse request body %+v", caller, err)
		d.ServeErrResponse(ctx, err)
		return
	}

	discount, err := d.discountService.Update(ctxReq, uuid.MustParse(id), req)
	if err != nil {
		log.Error(ctxReq, "[%v] failed to update discount %+v", caller, err)
		d.ServeErrResponse(ctx, err)
		return
	}
	d.ServeSuccessResponse(ctx, response.ToDiscountResponse(*discount))
}

func (d DiscountController) SetStatus(ctx *gin.Context) {
	caller := "DiscountController.SetStatus"
	ctxReq := ctx.Request.Context()

	id := ctx.Param("id")
	isActive := ctx.Query("is_active") == "true"

	discount, err := d.discountService.SetStatus(ctxReq, uuid.MustParse(id), isActive)
	if err != nil {
		log.Error(ctxReq, "[%v] failed to set discount status %+v", caller, err)
		d.ServeErrResponse(ctx, err)
		return
	}
	d.ServeSuccessResponse(ctx, response.ToDiscountResponse(*discount))
}
