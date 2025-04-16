package controller

import (
	"TTCS/src/common/log"
	"TTCS/src/core/service"
	"TTCS/src/present/httpui/response"
	"github.com/gin-gonic/gin"
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
	d.ServeSuccessResponse(ctx, response.ToDiscountResponse(discount))
	return
}
