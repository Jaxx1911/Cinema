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
		log.Error(ctxReq, "[%v] failed to create showtimes %+v", caller, err)
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
		log.Error(ctxReq, "[%v] failed to get showtimes %+v", caller, err)
		s.ServeErrResponse(ctx, err)
		return
	}
	s.ServeSuccessResponse(ctx, response.ToListShowtimeWithRoom(showtimes))
}

func (s *ShowtimeController) GetByCinemaId(ctx *gin.Context) {
	caller := "ShowtimeController.GetByCinemaId"
	ctxReq := ctx.Request.Context()

	var req request.GetShowtimesByCinemaIdFilter
	if err := ctx.ShouldBindQuery(&req); err != nil {
		log.Error(ctxReq, "[%v] invalid query %+v", caller, err)
		s.ServeErrResponse(ctx, err)
		return
	}

	showtimes, err := s.ShowtimeService.GetByCinemaFilter(ctxReq, req)
	if err != nil {
		log.Error(ctxReq, "[%v] failed to get showtimes %+v", caller, err)
		s.ServeErrResponse(ctx, err)
		return
	}
	s.ServeSuccessResponse(ctx, response.ToListShowtimeWithRoom(showtimes))
}

func (s *ShowtimeController) GetByRoomId(ctx *gin.Context) {
	caller := "ShowtimeController.GetByCinemaId"
	ctxReq := ctx.Request.Context()

	var req request.GetShowtimesByRoomIdFilter
	if err := ctx.ShouldBindQuery(&req); err != nil {
		log.Error(ctxReq, "[%v] invalid query %+v", caller, err)
		s.ServeErrResponse(ctx, err)
		return
	}

	showtimes, err := s.ShowtimeService.GetByRoomFilter(ctxReq, req)
	if err != nil {
		log.Error(ctxReq, "[%v] failed to get showtimes %+v", caller, err)
		s.ServeErrResponse(ctx, err)
		return
	}
	s.ServeSuccessResponse(ctx, response.ToListShowtimeWithMovieAndRoom(showtimes))
}

func (s *ShowtimeController) GetById(ctx *gin.Context) {
	caller := "ShowtimeController.GetById"
	ctxReq := ctx.Request.Context()

	id := ctx.Param("id")

	showtime, err := s.ShowtimeService.GetById(ctxReq, id)
	if err != nil {
		log.Error(ctxReq, "[%v] failed to get showtimes %+v", caller, err)
		s.ServeErrResponse(ctx, err)
		return
	}
	s.ServeSuccessResponse(ctx, response.ToShowtimeDetailResponse(showtime))
}

func (s *ShowtimeController) GetList(ctx *gin.Context) {
	ctxReq := ctx.Request.Context()
	caller := "ShowtimeController.GetList"

	var req request.GetListShowtime
	if err := ctx.ShouldBindQuery(&req); err != nil {
		log.Error(ctxReq, "[%v] invalid param %+v", caller, err)
		s.ServeErrResponse(ctx, err)
		return
	}

	// Set default values for pagination
	req.SetDefaults()

	showtimes, total, err := s.ShowtimeService.GetList(ctxReq, req)
	if err != nil {
		log.Error(ctxReq, "[%v] failed to get showtimes: %v", caller, err)
		s.ServeErrResponse(ctx, err)
		return
	}

	s.ServeSuccessResponse(ctx, response.MetaData{
		Data:       response.ToListShowtimeWithRoomAndSoldTickets(showtimes),
		TotalCount: total,
	})
}

func (s *ShowtimeController) Update(ctx *gin.Context) {
	caller := "ShowtimeController.Update"
	ctxReq := ctx.Request.Context()

	id := ctx.Param("id")

	var req request.UpdateShowtime
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error(ctxReq, "[%v] invalid param %+v", caller, err)
		s.ServeErrResponse(ctx, err)
		return
	}

	showtime, err := s.ShowtimeService.Update(ctxReq, id, req)
	if err != nil {
		log.Error(ctxReq, "[%v] failed to update showtime %+v", caller, err)
		s.ServeErrResponse(ctx, err)
		return
	}

	s.ServeSuccessResponse(ctx, response.ToShowtimeResponse(*showtime))
}

func (s *ShowtimeController) Delete(ctx *gin.Context) {
	caller := "ShowtimeController.Delete"
	ctxReq := ctx.Request.Context()

	id := ctx.Param("id")

	err := s.ShowtimeService.Delete(ctxReq, id)
	if err != nil {
		log.Error(ctxReq, "[%v] failed to delete showtime %+v", caller, err)
		s.ServeErrResponse(ctx, err)
		return
	}

	s.ServeSuccessResponse(ctx, gin.H{"message": "Showtime deleted successfully"})
}

func (s *ShowtimeController) CheckAvailability(ctx *gin.Context) {
	caller := "ShowtimeController.CheckAvailability"
	ctxReq := ctx.Request.Context()

	var req request.CheckShowtimeAvailability
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error(ctxReq, "[%v] invalid param %+v", caller, err)
		s.ServeErrResponse(ctx, err)
		return
	}

	availabilityResp, err := s.ShowtimeService.CheckShowtimeAvailability(ctxReq, req)
	if err != nil {
		log.Error(ctxReq, "[%v] failed to check showtime availability %+v", caller, err)
		s.ServeErrResponse(ctx, err)
		return
	}

	s.ServeSuccessResponse(ctx, response.ToShowtimeAvailabilityResponse(availabilityResp))
}

func (s *ShowtimeController) CheckShowtimesAvailability(ctx *gin.Context) {
	caller := "ShowtimeController.CheckShowtimesAvailability"
	ctxReq := ctx.Request.Context()

	var req request.CheckShowtimesAvailability
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error(ctxReq, "[%v] invalid param %+v", caller, err)
		s.ServeErrResponse(ctx, err)
		return
	}

	result, err := s.ShowtimeService.CheckShowtimesAvailability(ctxReq, req)
	if err != nil {
		log.Error(ctxReq, "[%v] failed to check showtimes availability %+v", caller, err)
		s.ServeErrResponse(ctx, err)
		return
	}

	s.ServeSuccessResponse(ctx, response.ToShowtimesAvailabilityResponse(result))
}

func (s *ShowtimeController) CreateShowtimes(ctx *gin.Context) {
	caller := "ShowtimeController.CreateShowtimes"
	ctxReq := ctx.Request.Context()

	var req request.CreateShowtimes
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error(ctxReq, "[%v] invalid param %+v", caller, err)
		s.ServeErrResponse(ctx, err)
		return
	}

	result, err := s.ShowtimeService.CreateShowtimes(ctxReq, req)
	if err != nil {
		log.Error(ctxReq, "[%v] failed to create showtimes %+v", caller, err)
		s.ServeErrResponse(ctx, err)
		return
	}

	s.ServeSuccessResponse(ctx, response.ToCreateShowtimesResponse(result))
}
