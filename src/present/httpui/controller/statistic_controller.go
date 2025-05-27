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

// GetMovieRevenue godoc
// @Summary Get movie revenue statistics
// @Description Get revenue statistics by movie with date range
// @Tags statistics
// @Accept json
// @Produce json
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} response.MovieRevenueResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /statistic/movie-revenue [get]
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

// GetCinemaRevenue godoc
// @Summary Get cinema revenue statistics
// @Description Get revenue statistics by cinema with date range
// @Tags statistics
// @Accept json
// @Produce json
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} response.CinemaRevenueResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /statistic/cinema-revenue [get]
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

// GetComboStatistics godoc
// @Summary Get combo statistics
// @Description Get combo sales statistics
// @Tags statistics
// @Accept json
// @Produce json
// @Success 200 {object} response.ComboStatisticResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /statistic/combo [get]
func (s *StatisticController) GetComboStatistics(ctx *gin.Context) {
	caller := "StatisticController.GetComboStatistics"
	ctxReq := ctx.Request.Context()

	result, err := s.statisticService.GetComboStatistics(ctxReq)
	if err != nil {
		log.Error(ctxReq, "[%v] failed to get combo statistics %+v", caller, err)
		s.ServeErrResponse(ctx, err)
		return
	}

	s.ServeSuccessResponse(ctx, result)
}
