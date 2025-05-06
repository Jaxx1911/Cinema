package controller

import (
	"TTCS/src/common/log"
	"TTCS/src/core/service"
	"TTCS/src/present/httpui/request"
	"TTCS/src/present/httpui/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		log.Error(ctxReq, "[%v] get cinema failed %+v", caller, err)
		c.ServeErrResponse(ctx, err)
		return
	}
	c.ServeSuccessResponse(ctx, response.ToListCinemaResponse(cinemas))
}

func (c *CinemaController) GetListByCity(ctx *gin.Context) {
	caller := "CinemaController.GetListByCity"
	ctxReq := ctx.Request.Context()

	var req request.GetCinemaRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		log.Error(ctxReq, "[%v] invalid param %+v", caller, err)
		c.ServeErrResponse(ctx, err)
		return
	}

	req.MappingCity()

	cinemas, err := c.cinemaService.GetListByCity(ctxReq, req.City)
	if err != nil {
		log.Error(ctxReq, "[%v] get cinema failed %+v", caller, err)
		c.ServeErrResponse(ctx, err)
		return
	}
	c.ServeSuccessResponse(ctx, response.ToListCinemaResponse(cinemas))
}

func (c *CinemaController) GetFacilities(ctx *gin.Context) {
	caller := "CinemaController.GetFacilities"
	ctxReq := ctx.Request.Context()

	var req request.GetCinemaRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		log.Error(ctxReq, "[%v] invalid param %+v", caller, err)
		c.ServeErrResponse(ctx, err)
		return
	}
	req.MappingCity()

	cinemasWithFacilities, err := c.cinemaService.GetWithRoomsByCity(ctxReq, req.City)
	if err != nil {
		log.Error(ctxReq, "[%v] get cinema failed %+v", caller, err)
		c.ServeErrResponse(ctx, err)
		return
	}
	c.ServeSuccessResponse(ctx, response.ToListCinemaWithFacilitiesResponse(cinemasWithFacilities))
}

func (c *CinemaController) GetCinemaDetail(ctx *gin.Context) {
	caller := "CinemaController.GetCinema"
	ctxReq := ctx.Request.Context()

	id := ctx.Param("id")

	cinema, err := c.cinemaService.GetDetail(ctxReq, uuid.MustParse(id))
	if err != nil {
		log.Error(ctxReq, "[%v] get cinema failed", caller, err)
		c.ServeErrResponse(ctx, err)
		return
	}
	c.ServeSuccessResponse(ctx, response.ToCinemaWithFacilitiesResponse(cinema))
}

func (c *CinemaController) Create(ctx *gin.Context) {
	caller := "CinemaController.Create"
	ctxReq := ctx.Request.Context()

	var req request.CreateCinemaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error(ctxReq, "[%v] invalid param %+v", caller, err)
		c.ServeErrResponse(ctx, err)
		return
	}

	cinema, err := c.cinemaService.Create(ctx, req)
	if err != nil {
		log.Error(ctxReq, "[%v] create cinema failed", caller, err)
		c.ServeErrResponse(ctx, err)
		return
	}
	c.ServeSuccessResponse(ctx, response.ToCinemaResponse(cinema))
}

func (c *CinemaController) Update(ctx *gin.Context) {
	caller := "CinemaController.Update"
	ctxReq := ctx.Request.Context()

	var req request.UpdateCinemaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error(ctxReq, "[%v] invalid param %+v", caller, err)
		c.ServeErrResponse(ctx, err)
		return
	}

	cinema, err := c.cinemaService.Update(ctxReq, req)
	if err != nil {
		log.Error(ctxReq, "[%v] update cinema failed", caller, err)
		c.ServeErrResponse(ctx, err)
		return
	}
	c.ServeSuccessResponse(ctx, response.ToCinemaResponse(cinema))
}

func (c *CinemaController) Delete(ctx *gin.Context) {
	caller := "CinemaController.Delete"
	ctxReq := ctx.Request.Context()

	id := ctx.Param("id")
	if err := c.cinemaService.Delete(ctx, uuid.MustParse(id)); err != nil {
		log.Error(ctxReq, "[%v] delete cinema failed", caller, err)
		c.ServeErrResponse(ctx, err)
		return
	}
	c.ServeSuccessResponse(ctx, true)
}
