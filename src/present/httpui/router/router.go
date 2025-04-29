package router

import (
	"TTCS/src/common/configs"
	"TTCS/src/present/httpui/controller"
	"TTCS/src/present/httpui/middleware"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	cors "github.com/rs/cors/wrapper/gin"
	"go.uber.org/fx"
)

type IRouter struct {
	fx.In
	Engine              *gin.Engine
	AuthHolder          *middleware.AuthMiddleware
	AuthController      *controller.AuthController
	UserController      *controller.UserController
	MovieController     *controller.MovieController
	ShowtimeController  *controller.ShowtimeController
	CinemaController    *controller.CinemaController
	SeatController      *controller.SeatController
	ComboController     *controller.ComboController
	DiscountController  *controller.DiscountController
	OrderController     *controller.OrderController
	PaymentController   *controller.PaymentController
	WebSocketController *controller.WebSocketController
}

func RegisterHandler(engine *gin.Engine) {
	engine.Use(middleware.Log())
}

func RegisterGinRouters(in IRouter) {
	in.Engine.Use(cors.AllowAll())

	group := in.Engine.Group(configs.GetConfig().Server.Prefix)
	group.GET("/ping")

	registerRouters(group, in)

	registerWebSocket(in.Engine, in)
}

func registerRouters(r *gin.RouterGroup, in IRouter) {
	registerAuthRouters(r, in)
	registerUsersRouters(r, in)
	registerMovieRouters(r, in)
	registerShowtimeRouter(r, in)
	registerCinemaRouter(r, in)
	registerSeatRouter(r, in)
	registerComboRouter(r, in)
	registerDiscountRouter(r, in)
	registerOrderRouter(r, in)
	registerPaymentRouter(r, in)
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
		showtimeRouter.GET("/:id", in.ShowtimeController.GetById)
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

func registerSeatRouter(root *gin.RouterGroup, in IRouter) {
	seatRouter := root.Group("/seat")
	{
		seatRouter.GET("/room/:id", in.SeatController.GetByRoomId)
	}
}

func registerComboRouter(root *gin.RouterGroup, in IRouter) {
	comboRouter := root.Group("/combo")
	{
		comboRouter.GET("", in.ComboController.GetList)
	}
}

func registerOrderRouter(root *gin.RouterGroup, in IRouter) {
	orderRouter := root.Group("/order")
	{
		orderRouter.Use(in.AuthHolder.RequireAuth())
		orderRouter.POST("", in.OrderController.CreateOrder)
		orderRouter.GET("/qr/:id", in.OrderController.GetOrderDetailsWithQr)
		orderRouter.GET("/:id", in.OrderController.GetOrderDetails)
	}
}

func registerPaymentRouter(root *gin.RouterGroup, in IRouter) {
	paymentRouter := root.Group("/payment")
	{
		paymentRouter.POST("/callback", in.PaymentController.CallBack)
		paymentRouter.GET("", in.PaymentController.GetListByUserId)
	}

}

func registerDiscountRouter(root *gin.RouterGroup, in IRouter) {
	discountRouter := root.Group("/discount")
	{
		discountRouter.GET("", in.DiscountController.GetDiscountByCode)
	}
}

func registerWebSocket(root *gin.Engine, in IRouter) {
	websocketGroup := root.Group("/ws", func(c *gin.Context) {
		if websocket.IsWebSocketUpgrade(c.Request) {
			c.Set("allowed", true)
			c.Next()
		}
		return
	})
	websocketGroup.GET("", in.WebSocketController.HandleWebSocket)
}
