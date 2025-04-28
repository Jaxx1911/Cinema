package controller

import (
	"TTCS/src/common/log"
	"TTCS/src/core/service"
	"TTCS/src/present/httpui/response"
	"github.com/gin-gonic/gin"
)

type ComboController struct {
	*BaseController
	comboService service.ComboService
}

func NewComboController(baseController *BaseController, comboService service.ComboService) *ComboController {
	return &ComboController{
		baseController,
		comboService,
	}
}

func (c *ComboController) GetList(ctx *gin.Context) {
	caller := "ComboController.GetList"
	ctxReq := ctx.Request.Context()

	combos, err := c.comboService.GetList(ctxReq)
	if err != nil {
		log.Error(ctx, "[%v] get combos %+v", caller, err)
		c.ServeErrResponse(ctx, err)
		return
	}
	c.ServeSuccessResponse(ctx, response.ToListComboResponse(combos))
}
