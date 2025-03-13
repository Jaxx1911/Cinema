package router

import (
	"TTCS/src/common/configs"
	"TTCS/src/present/httpui/controller"
	"TTCS/src/present/httpui/middleware"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"go.uber.org/fx"
)

type IRouter struct {
	fx.In
	Engine         *gin.Engine
	AuthHolder     *middleware.AuthMiddleware
	AuthController *controller.AuthController
	UserController *controller.UserController
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
	registerAuthRouters(r, in)
	registerUsersRouters(r, in)
}
func registerAuthRouters(root *gin.RouterGroup, in IRouter) {
	authRouter := root.Group("/auth")
	{
		authRouter.GET("/otp/:email", in.AuthController.GetOTP)
		authRouter.POST("/signup", in.AuthController.SignUp)
		authRouter.POST("/login", in.AuthController.Login)
		authRouter.Use(in.AuthHolder.RequireAuth())
		authRouter.POST("/change-password", in.AuthController.ChangePassword)
		//authRouter.GET("/login-google", in.AuthController.LoginGoogle)
		//authRouter.GET("/callback-by-google", in.AuthController.CallbackGoogle)
	}
}
func registerUsersRouters(root *gin.RouterGroup, in IRouter) {

	userRouter := root.Group("/user")
	{
		userRouter.Use(in.AuthHolder.RequireAuth())
		userRouter.PUT("", in.UserController.UpdateInfo)
		userRouter.POST("", in.UserController.Create)
		userRouter.GET("", in.UserController.GetList)
		userRouter.GET("/:id", in.UserController.GetDetail)
		userRouter.GET("/payments", in.UserController.GetPayments)
		userRouter.GET("/orders", in.UserController.GetOrders)
		userRouter.PUT("/avatar", in.UserController.ChangeAvatar)
	}
}
