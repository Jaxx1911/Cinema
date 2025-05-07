package controller

import (
	"TTCS/src/common/log"
	"TTCS/src/core/service"
	"TTCS/src/present/httpui/request"
	"TTCS/src/present/httpui/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MovieController struct {
	*BaseController
	movieService *service.MovieService
}

func NewMovieController(baseController *BaseController, movieService *service.MovieService) *MovieController {
	return &MovieController{
		BaseController: baseController,
		movieService:   movieService,
	}
}

func (m *MovieController) GetList(ctx *gin.Context) {
	ctxReq := ctx.Request.Context()
	caller := "UserController.GetList"

	var page request.GetListMovie
	if err := ctx.ShouldBindQuery(&page); err != nil {
		log.Error(ctxReq, "[%v] invalid param %+v", caller, err)
		m.ServeErrResponse(ctx, err)
		return
	}

	page.Page.SetDefaults()

	movies, total, err := m.movieService.GetList(ctxReq, page)
	if err != nil {
		log.Error(ctxReq, "[%v] get movie list failed", caller, err)
		m.ServeErrResponse(ctx, err)
		return
	}
	m.ServeSuccessResponse(ctx, response.MetaData{
		Data:       response.ToListMoviesResponse(movies),
		TotalCount: total,
	})
}

func (m *MovieController) GetListByStatus(ctx *gin.Context) {
	ctxReq := ctx.Request.Context()
	caller := "UserController.GetListByStatus"

	status := ctx.Query("status")

	var page request.Page
	if err := ctx.ShouldBindQuery(&page); err != nil {
		log.Error(ctxReq, "[%v] invalid param %+v", caller, err)
		m.ServeErrResponse(ctx, err)
		return
	}

	page.SetDefaults()

	movies, err := m.movieService.GetListByStatus(ctxReq, page, status)
	if err != nil {
		log.Error(ctxReq, "[%v] get movie list failed", caller, err)
		m.ServeErrResponse(ctx, err)
		return
	}
	m.ServeSuccessResponse(ctx, response.ToListMoviesResponse(movies))
}

func (m *MovieController) GetDetail(ctx *gin.Context) {
	ctxReq := ctx.Request.Context()
	caller := "UserController.GetDetail"

	movie, err := m.movieService.GetDetail(ctxReq, uuid.MustParse(ctx.Param("id")))
	if err != nil {
		log.Error(ctxReq, "[%v] get movie detail failed", caller, err)
		m.ServeErrResponse(ctx, err)
		return
	}
	m.ServeSuccessResponse(ctx, response.ToMovieDetailResponse(movie))
}

func (m *MovieController) Create(ctx *gin.Context) {
	ctxReq := ctx.Request.Context()
	caller := "UserController.Create"

	var req request.CreateMovieRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.Error(ctxReq, "[%v] invalid param %+v", caller, err)
		m.ServeErrResponse(ctx, err)
		return
	}

	movie, err := m.movieService.Create(ctxReq, req)
	if err != nil {
		log.Error(ctxReq, "[%v] create movie failed %+v", caller, err)
		m.ServeErrResponse(ctx, err)
		return
	}

	m.ServeSuccessResponse(ctx, movie)
}

func (m *MovieController) Update(ctx *gin.Context) {
	ctxReq := ctx.Request.Context()
	caller := "UserController.Update"

	var req request.UpdateMovieRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.Error(ctxReq, "[%v] invalid param %+v", caller, err)
		m.ServeErrResponse(ctx, err)
		return
	}
	id := ctx.Param("id")
	req.Id = uuid.MustParse(id)

	movie, err := m.movieService.Update(ctxReq, req)
	if err != nil {
		log.Error(ctxReq, "[%v] update movie failed %+v", caller, err)
		m.ServeErrResponse(ctx, err)
		return
	}
	m.ServeSuccessResponse(ctx, response.ToMovieDetailResponse(movie))
}

func (m *MovieController) GetListInDateRange(ctx *gin.Context) {
	ctxReq := ctx.Request.Context()
	caller := "UserController.GetListInDateRange"

	var page request.Page
	if err := ctx.ShouldBindQuery(&page); err != nil {
		log.Error(ctxReq, "[%v] invalid param %+v", caller, err)
		m.ServeErrResponse(ctx, err)
		return
	}

	page.SetDefaults()

	movies, err := m.movieService.GetListInDateRange(ctxReq)
	if err != nil {
		log.Error(ctxReq, "[%v] get movie list failed", caller, err)
		m.ServeErrResponse(ctx, err)
		return
	}
	m.ServeSuccessResponse(ctx, response.ToListMoviesResponse(movies))
}

func (m *MovieController) StopMovie(ctx *gin.Context) {
	ctxReq := ctx.Request.Context()
	caller := "UserController.GetListOutDateRange"

	id := ctx.Param("id")

	if err := m.movieService.StopMovie(ctxReq, uuid.MustParse(id)); err != nil {
		log.Error(ctxReq, "[%v] stop movie failed %+v", caller, err)
		m.ServeErrResponse(ctx, err)
		return
	}
	m.ServeSuccessResponse(ctx, true)
}

func (m *MovieController) ReshowMovie(ctx *gin.Context) {
	ctxReq := ctx.Request.Context()
	caller := "UserController.GetListOutDateRange"

	id := ctx.Param("id")

	if err := m.movieService.ReshowMovie(ctxReq, uuid.MustParse(id)); err != nil {
		log.Error(ctxReq, "[%v] stop movie failed %+v", caller, err)
		m.ServeErrResponse(ctx, err)
		return
	}
	m.ServeSuccessResponse(ctx, true)
}
