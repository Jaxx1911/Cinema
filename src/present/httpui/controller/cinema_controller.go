package controller

import (
	"TTCS/src/common/log"
	"TTCS/src/core/service"
	"TTCS/src/present/httpui/response"
	"github.com/gin-gonic/gin"
)

type CinemaController struct {
	*BaseController
	cinemaService *service.CinemaService
}

func NewCinemaController(baseController *BaseController, cinemaService *service.CinemaService) *CinemaController {
	return &CinemaController{
		BaseController: baseController,
		cinemaService:  cinemaService,
	}
}

func (c *CinemaController) GetList(ctx *gin.Context) {
	caller := "CinemaController.GetList"
	ctxReq := ctx.Request.Context()

	cinemas, err := c.cinemaService.GetList(ctxReq)
	if err != nil {
		log.Error(ctxReq, "[%v] get cinema failed", caller, err)
		c.ServeErrResponse(ctx, err)
		return
	}
	c.ServeSuccessResponse(ctx, response.ToListCinemaResponse(cinemas))
}
