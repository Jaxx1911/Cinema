package controller

import (
	"TTCS/src/common/log"
	"TTCS/src/core/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RoomController struct {
	*BaseController
	roomService *service.RoomService
}

func NewRoomController(baseController *BaseController, roomService *service.RoomService) *RoomController {
	return &RoomController{
		BaseController: baseController,
		roomService:    roomService,
	}
}

func (r *RoomController) Deactivate(ctx *gin.Context) {
	caller := "RoomController.Deactive"
	ctxReq := ctx.Request.Context()

	isActiveString := ctx.Query("is_active")
	id := ctx.Param("id")

	var isActive bool
	if isActiveString == "true" {
		isActive = true
	} else {
		isActive = false
	}
	if err := r.roomService.Deactivate(ctx, uuid.MustParse(id), isActive); err != nil {
		log.Error(ctxReq, "[%v] failed to deactivate +%v", caller, err)
		r.ServeErrResponse(ctx, err)
	}
	r.ServeSuccessResponse(ctx, true)
}
