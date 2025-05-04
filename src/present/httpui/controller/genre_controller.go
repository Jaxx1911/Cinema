package controller

import (
	"TTCS/src/common/log"
	"TTCS/src/core/service"
	"TTCS/src/present/httpui/response"
	"github.com/gin-gonic/gin"
)

type GenreController struct {
	*BaseController
	genService *service.GenreService
}

func NewGenreController(baseController *BaseController, genService *service.GenreService) *GenreController {
	return &GenreController{
		BaseController: baseController,
		genService:     genService,
	}
}

func (g *GenreController) GetGenres(c *gin.Context) {
	caller := "GenreController.GetGenres"
	ctxReq := c.Request.Context()

	genres, err := g.genService.GetGenres(ctxReq)
	if err != nil {
		log.Error(ctxReq, "[%v] get genre list failed", caller, err)
		g.ServeErrResponse(c, err)
		return
	}
	g.ServeSuccessResponse(c, response.ToGenresResponse(genres))
	return
}
