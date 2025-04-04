package controller

import (
	"TTCS/src/common/log"
	"TTCS/src/core/service"
	"TTCS/src/present/httpui/response"
	"github.com/gin-gonic/gin"
)

type SeatController struct {
	*BaseController
	seatService *service.SeatService
}

func NewSeatController(baseController *BaseController, seatService *service.SeatService) *SeatController {
	return &SeatController{
		seatService:    seatService,
		BaseController: baseController,
	}
}

func (s SeatController) GetByRoomId(ctx *gin.Context) {
	caller := "SeatController.GetByRoomId"
	ctxReq := ctx.Request.Context()

	id := ctx.Param("id")

	seats, err := s.seatService.GetByRoomId(ctxReq, id)
	if err != nil {
		log.Error(ctxReq, "[%v] failed to get seats %+v", caller, err)
		s.ServeErrResponse(ctx, err)
	}
	s.ServeSuccessResponse(ctx, response.ToListSeatResponse(seats))
}
