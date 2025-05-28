package controller

import (
	"TTCS/src/common/log"
	"TTCS/src/core/service"
	"TTCS/src/present/httpui/request"

	"github.com/gin-gonic/gin"
)

type StatisticController struct {
	*BaseController
	statisticService *service.StatisticService
}

func NewStatisticController(baseController *BaseController, statisticService *service.StatisticService) *StatisticController {
	return &StatisticController{
		BaseController:   baseController,
		statisticService: statisticService,
	}
}

func (s *StatisticController) GetMovieRevenue(ctx *gin.Context) {
	caller := "StatisticController.GetMovieRevenue"
	ctxReq := ctx.Request.Context()

	var req request.StatisticDateRange
	if err := ctx.ShouldBindQuery(&req); err != nil {
		log.Error(ctxReq, "[%v] invalid request %+v", caller, err)
		s.ServeErrResponse(ctx, err)
		return
	}

	result, err := s.statisticService.GetMovieRevenue(ctxReq, req)
	if err != nil {
		log.Error(ctxReq, "[%v] failed to get movie revenue %+v", caller, err)
		s.ServeErrResponse(ctx, err)
		return
	}

	s.ServeSuccessResponse(ctx, result)
}

func (s *StatisticController) GetCinemaRevenue(ctx *gin.Context) {
	caller := "StatisticController.GetCinemaRevenue"
	ctxReq := ctx.Request.Context()

	var req request.StatisticDateRange
	if err := ctx.ShouldBindQuery(&req); err != nil {
		log.Error(ctxReq, "[%v] invalid request %+v", caller, err)
		s.ServeErrResponse(ctx, err)
		return
	}

	result, err := s.statisticService.GetCinemaRevenue(ctxReq, req)
	if err != nil {
		log.Error(ctxReq, "[%v] failed to get cinema revenue %+v", caller, err)
		s.ServeErrResponse(ctx, err)
		return
	}

	s.ServeSuccessResponse(ctx, result)
}

func (s *StatisticController) GetComboStatistics(ctx *gin.Context) {
	caller := "StatisticController.GetComboStatistics"
	ctxReq := ctx.Request.Context()

	var req request.StatisticDateRange
	if err := ctx.ShouldBindQuery(&req); err != nil {
		log.Error(ctxReq, "[%v] invalid request %+v", caller, err)
		s.ServeErrResponse(ctx, err)
		return
	}

	result, err := s.statisticService.GetComboStatistics(ctxReq, req)
	if err != nil {
		log.Error(ctxReq, "[%v] failed to get combo statistics %+v", caller, err)
		s.ServeErrResponse(ctx, err)
		return
	}

	s.ServeSuccessResponse(ctx, result)
}
