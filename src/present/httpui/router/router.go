package router

import (
	"TTCS/src/common/configs"
	"TTCS/src/present/httpui/middleware"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"go.uber.org/fx"
)

type IRouter struct {
	fx.In
	Engine *gin.Engine
	//controller
}

func RegisterHandler(engine *gin.Engine) {
	engine.Use(middleware.Log())
}

func RegisterGinRouters(in IRouter) {
	in.Engine.Use(cors.AllowAll())

	group := in.Engine.Group(configs.GetConfig().Server.Prefix)
	group.GET("/ping")

	registerRouters(group, in)
}

func registerRouters(r *gin.RouterGroup, in IRouter) {
}
