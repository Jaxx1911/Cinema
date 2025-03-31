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
	Engine             *gin.Engine
	AuthHolder         *middleware.AuthMiddleware
	AuthController     *controller.AuthController
	UserController     *controller.UserController
	MovieController    *controller.MovieController
	ShowtimeController *controller.ShowtimeController
	CinemaController   *controller.CinemaController
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
	registerMovieRouters(r, in)
	registerShowtimeRouter(r, in)
	registerCinemaRouter(r, in)
}
func registerAuthRouters(root *gin.RouterGroup, in IRouter) {
	authRouter := root.Group("/auth")
	{
		authRouter.GET("/otp/:email", in.AuthController.SignUpOTP)
		authRouter.POST("/signup", in.AuthController.SignUp)
		authRouter.POST("/login", in.AuthController.Login)
		authRouter.GET("/reset-otp/:email", in.AuthController.ResetOTP)
		authRouter.POST("/reset-password", in.AuthController.ResetPassword)
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
		userRouter.GET("/detail", in.UserController.GetDetail)
		userRouter.GET("/payments", in.UserController.GetPayments)
		userRouter.GET("/orders", in.UserController.GetOrders)
		userRouter.PUT("/avatar", in.UserController.ChangeAvatar)
	}
}

func registerMovieRouters(root *gin.RouterGroup, in IRouter) {
	movieRouter := root.Group("/movie")
	{
		movieRouter.GET("", in.MovieController.GetList)
		movieRouter.GET("/range", in.MovieController.GetListInDateRange)
		movieRouter.GET("/:id", in.MovieController.GetDetail)
		movieRouter.Use(in.AuthHolder.RequireAuth())
		movieRouter.POST("", in.MovieController.Create)
		movieRouter.PUT("/:id", in.MovieController.Update)
		movieRouter.PUT("/:id/poster", in.MovieController.UpdatePoster)
	}
}

func registerShowtimeRouter(root *gin.RouterGroup, in IRouter) {
	showtimeRouter := root.Group("/showtime")
	{
		showtimeRouter.POST("", in.ShowtimeController.Create)
		showtimeRouter.GET("", in.ShowtimeController.GetByUserFilter)
		showtimeRouter.GET("/cinema", in.ShowtimeController.GetByCinemaId)
	}
}

func registerCinemaRouter(root *gin.RouterGroup, in IRouter) {
	cinemaRouter := root.Group("/cinema")
	{
		cinemaRouter.GET("", in.CinemaController.GetList)
		cinemaRouter.GET("/facilities", in.CinemaController.GetFacilities)
		cinemaRouter.GET("/:id", in.CinemaController.GetCinemaDetail)
	}
}
