package controller

import (
	"TTCS/src/common/log"
	"TTCS/src/core/service"
	"TTCS/src/present/httpui/request"
	"TTCS/src/present/httpui/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ComboController struct {
	*BaseController
	comboService service.ComboService
}

func NewComboController(baseController *BaseController, comboService service.ComboService) *ComboController {
	return &ComboController{
		BaseController: baseController,
		comboService:   comboService,
	}
}

func (c *ComboController) GetList(ctx *gin.Context) {
	caller := "ComboController.GetList"
	ctxReq := ctx.Request.Context()

	combos, err := c.comboService.GetList(ctxReq)
	if err != nil {
		log.Error(ctxReq, "[%v] get combos %+v", caller, err)
		c.ServeErrResponse(ctx, err)
		return
	}
	c.ServeSuccessResponse(ctx, response.ToListComboResponse(combos))
}

func (c *ComboController) GetDetail(ctx *gin.Context) {
	caller := "ComboController.GetDetail"
	ctxReq := ctx.Request.Context()

	id := ctx.Param("id")
	combo, err := c.comboService.GetDetail(ctxReq, uuid.MustParse(id))
	if err != nil {
		log.Error(ctxReq, "[%v] get combo detail failed %+v", caller, err)
		c.ServeErrResponse(ctx, err)
		return
	}
	c.ServeSuccessResponse(ctx, response.ToComboResponse(combo))
}

func (c *ComboController) Create(ctx *gin.Context) {
	caller := "ComboController.Create"
	ctxReq := ctx.Request.Context()

	var req request.CreateComboRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.Error(ctxReq, "[%v] invalid param %+v", caller, err)
		c.ServeErrResponse(ctx, err)
		return
	}

	combo, err := c.comboService.Create(ctxReq, req)
	if err != nil {
		log.Error(ctxReq, "[%v] create combo failed %+v", caller, err)
		c.ServeErrResponse(ctx, err)
		return
	}
	c.ServeSuccessResponse(ctx, response.ToComboResponse(combo))
}

func (c *ComboController) Update(ctx *gin.Context) {
	caller := "ComboController.Update"
	ctxReq := ctx.Request.Context()

	id := ctx.Param("id")
	var req request.UpdateComboRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.Error(ctxReq, "[%v] invalid param %+v", caller, err)
		c.ServeErrResponse(ctx, err)
		return
	}

	combo, err := c.comboService.Update(ctxReq, uuid.MustParse(id), req)
	if err != nil {
		log.Error(ctxReq, "[%v] update combo failed %+v", caller, err)
		c.ServeErrResponse(ctx, err)
		return
	}
	c.ServeSuccessResponse(ctx, response.ToComboResponse(combo))
}

func (c *ComboController) Delete(ctx *gin.Context) {
	caller := "ComboController.Delete"
	ctxReq := ctx.Request.Context()

	id := ctx.Param("id")
	err := c.comboService.Delete(ctxReq, uuid.MustParse(id))
	if err != nil {
		log.Error(ctxReq, "[%v] delete combo failed %+v", caller, err)
		c.ServeErrResponse(ctx, err)
		return
	}
	c.ServeSuccessResponse(ctx, true)
}
