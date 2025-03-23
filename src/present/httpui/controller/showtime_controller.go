package controller

import (
	"TTCS/src/common/log"
	"TTCS/src/core/service"
	"TTCS/src/present/httpui/request"
	"TTCS/src/present/httpui/response"
	"github.com/gin-gonic/gin"
)

type ShowtimeController struct {
	*BaseController
	ShowtimeService *service.ShowtimeService
}

func NewShowtimeController(baseController *BaseController, showtimeService *service.ShowtimeService) *ShowtimeController {
	return &ShowtimeController{
		BaseController:  baseController,
		ShowtimeService: showtimeService,
	}
}

func (s *ShowtimeController) Create(ctx *gin.Context) {
	caller := "ShowtimeController.Create"
	ctxReq := ctx.Request.Context()

	var req request.CreateShowtime
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error(ctxReq, "[%v] invalid param %+v", caller, err)
		s.ServeErrResponse(ctx, err)
		return
	}

	showtime, err := s.ShowtimeService.Create(ctxReq, req)
	if err != nil {
		log.Error(ctxReq, "[%v] invalid param %+v", caller, err)
		s.ServeErrResponse(ctx, err)
		return
	}
	s.ServeSuccessResponse(ctx, response.ToShowtimeResponse(*showtime))
}

func (s *ShowtimeController) GetByUserFilter(ctx *gin.Context) {
	caller := "ShowtimeController.GetByMovieId"
	ctxReq := ctx.Request.Context()

	var req request.GetShowtimesByUserFilter
	if err := ctx.ShouldBindQuery(&req); err != nil {
		log.Error(ctxReq, "[%v] invalid query %+v", caller, err)
		s.ServeErrResponse(ctx, err)
		return
	}

	showtimes, err := s.ShowtimeService.GetByUserFilter(ctxReq, req)
	if err != nil {
		log.Error(ctxReq, "[%v] invalid param %+v", caller, err)
		s.ServeErrResponse(ctx, err)
		return
	}
	s.ServeSuccessResponse(ctx, response.ToListShowtimeWithRoom(showtimes))
}
