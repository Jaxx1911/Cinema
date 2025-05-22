package controller

import (
	"TTCS/src/common/log"
	"TTCS/src/core/service"
	"TTCS/src/present/httpui/request"
	"TTCS/src/present/httpui/response"

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

func (r *RoomController) Create(ctx *gin.Context) {
	caller := "RoomController.Create"
	ctxReq := ctx.Request.Context()

	var req request.CreateRoomReq
	if err := ctx.ShouldBind(&req); err != nil {
		log.Error(ctxReq, "[%v] invalid request +%v", caller, err)
		r.ServeErrResponse(ctx, err)
		return
	}
	room, err := r.roomService.Create(ctx, req)
	if err != nil {
		log.Error(ctxReq, "[%v] room create failed +%v", caller, err)
		r.ServeErrResponse(ctx, err)
		return
	}
	r.ServeSuccessResponse(ctx, response.ToRoomResponse(room))
}

func (r *RoomController) Update(ctx *gin.Context) {
	caller := "RoomController.Update"
	ctxReq := ctx.Request.Context()

	id := ctx.Param("id")
	var req request.UpdateRoomReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error(ctxReq, "[%v] invalid request +%v", caller, err)
		r.ServeErrResponse(ctx, err)
		return
	}

	room, err := r.roomService.Update(ctxReq, uuid.MustParse(id), req)
	if err != nil {
		log.Error(ctxReq, "[%v] room update failed +%v", caller, err)
		r.ServeErrResponse(ctx, err)
		return
	}
	r.ServeSuccessResponse(ctx, response.ToRoomResponse(room))
}

func (r *RoomController) GetRoomById(ctx *gin.Context) {
	caller := "RoomController.GetRoomById"
	ctxReq := ctx.Request.Context()

	id := ctx.Param("id")

	room, err := r.roomService.GetRoomById(ctxReq, uuid.MustParse(id))
	if err != nil {
		log.Error(ctxReq, "[%v] get room failed +%v", caller, err)
		r.ServeErrResponse(ctx, err)
		return
	}
	r.ServeSuccessResponse(ctx, response.ToRoomResponse(room))
}

func (r *RoomController) GetList(ctx *gin.Context) {
	caller := "RoomController.GetList"
	ctxReq := ctx.Request.Context()

	cinemaId := ctx.Param("id")

	rooms, err := r.roomService.GetListByCinemaId(ctxReq, uuid.MustParse(cinemaId))
	if err != nil {
		log.Error(ctxReq, "[%v] failed to get rooms %+v", caller, err)
		r.ServeErrResponse(ctx, err)
		return
	}

	r.ServeSuccessResponse(ctx, response.ToListRoomResponse(rooms))
}
