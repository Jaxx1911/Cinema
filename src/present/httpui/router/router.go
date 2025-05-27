package router

import (
	"TTCS/src/present/httpui/controller"
	"TTCS/src/present/httpui/middleware"
	"os"

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
	RoomController      *controller.RoomController
	SeatController      *controller.SeatController
	ComboController     *controller.ComboController
	DiscountController  *controller.DiscountController
	OrderController     *controller.OrderController
	PaymentController   *controller.PaymentController
	GenreController     *controller.GenreController
	WebSocketController *controller.WebSocketController
	StatisticController *controller.StatisticController
}

func RegisterHandler(engine *gin.Engine) {
	engine.Use(middleware.Log())
}

func RegisterGinRouters(in IRouter) {
	in.Engine.Use(cors.AllowAll())

	group := in.Engine.Group(os.Getenv("SERVER_PREFIX"))
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
	registerGenreRouter(r, in)
	registerRoomRouter(r, in)
	registerStatisticRouter(r, in)
}
func registerAuthRouters(root *gin.RouterGroup, in IRouter) {
	authRouter := root.Group("/auth")
	{
		authRouter.GET("/otp/:email", in.AuthController.SignUpOTP)
		authRouter.POST("/signup", in.AuthController.SignUp)
		authRouter.POST("/login", in.AuthController.Login)
		authRouter.POST("/login/admin", in.AuthController.LoginAdmin)
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
		userRouter.GET("/me", in.UserController.GetMe)
		userRouter.PUT("/:id", in.UserController.Update)
		userRouter.PUT("", in.UserController.UpdateInfo)
		userRouter.POST("", in.UserController.Create)
		userRouter.GET("", in.UserController.GetList)
		userRouter.GET("/detail", in.UserController.GetDetail)
		userRouter.GET("/payments", in.UserController.GetPayments)
		userRouter.GET("/orders", in.UserController.GetOrders)
		userRouter.PUT("/avatar", in.UserController.ChangeAvatar)
		userRouter.DELETE("/:id", in.UserController.Delete)
	}
}

func registerMovieRouters(root *gin.RouterGroup, in IRouter) {
	movieRouter := root.Group("/movie")
	{
		movieRouter.GET("", in.MovieController.GetListByStatus)
		movieRouter.GET("/list", in.MovieController.GetList)
		movieRouter.GET("/range", in.MovieController.GetListInDateRange)
		movieRouter.GET("/:id", in.MovieController.GetDetail)
		//movieRouter.Use(in.AuthHolder.RequireAuth())
		movieRouter.POST("", in.MovieController.Create)
		movieRouter.PUT("/stop/:id", in.MovieController.StopMovie)
		movieRouter.PUT("/reshow/:id", in.MovieController.ReshowMovie)
		movieRouter.PUT("/:id", in.MovieController.Update)
	}
}

func registerShowtimeRouter(root *gin.RouterGroup, in IRouter) {
	showtimeRouter := root.Group("/showtime")
	{
		showtimeRouter.POST("", in.ShowtimeController.Create)
		showtimeRouter.POST("/batch", in.ShowtimeController.CreateShowtimes)
		showtimeRouter.POST("/check-availability", in.ShowtimeController.CheckAvailability)
		showtimeRouter.POST("/check-availabilities", in.ShowtimeController.CheckShowtimesAvailability)
		showtimeRouter.GET("", in.ShowtimeController.GetByUserFilter)
		showtimeRouter.GET("/cinema", in.ShowtimeController.GetByCinemaId)
		showtimeRouter.GET("/:id", in.ShowtimeController.GetById)
		showtimeRouter.GET("/room", in.ShowtimeController.GetByRoomId)
		//showtimeRouter.Use(in.AuthHolder.RequireAuth())
		showtimeRouter.GET("/list", in.ShowtimeController.GetList)
		showtimeRouter.PUT("/:id", in.ShowtimeController.Update)
		showtimeRouter.DELETE("/:id", in.ShowtimeController.Delete)
	}
}

func registerCinemaRouter(root *gin.RouterGroup, in IRouter) {
	cinemaRouter := root.Group("/cinema")
	{
		cinemaRouter.POST("", in.CinemaController.Create)
		cinemaRouter.PUT("/:id", in.CinemaController.Update)
		cinemaRouter.GET("", in.CinemaController.GetListByCity)
		cinemaRouter.GET("/list", in.CinemaController.GetList)
		cinemaRouter.GET("/facilities", in.CinemaController.GetFacilities)
		cinemaRouter.GET("/:id", in.CinemaController.GetCinemaDetail)
	}
}

func registerRoomRouter(root *gin.RouterGroup, in IRouter) {
	roomRouter := root.Group("/room")
	{
		roomRouter.POST("", in.RoomController.Create)
		roomRouter.GET("/:id", in.RoomController.GetRoomById)
		roomRouter.PUT("/:id", in.RoomController.Update)
		roomRouter.DELETE("/:id", in.RoomController.Delete)
		roomRouter.GET("/cinema/:id", in.RoomController.GetList)
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
		comboRouter.GET("/:id", in.ComboController.GetDetail)
		comboRouter.POST("", in.ComboController.Create)
		comboRouter.PUT("/:id", in.ComboController.Update)
		comboRouter.DELETE("/:id", in.ComboController.Delete)
	}
}

func registerOrderRouter(root *gin.RouterGroup, in IRouter) {
	orderRouter := root.Group("/order")
	{
		orderRouter.Use(in.AuthHolder.RequireAuth())
		orderRouter.POST("", in.OrderController.CreateOrder)
		orderRouter.GET("/qr/:id", in.OrderController.GetOrderDetailsWithQr)
		orderRouter.GET("/:id", in.OrderController.GetOrderDetails)
		orderRouter.DELETE("/:id", in.OrderController.DeleteOrder)
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
		discountRouter.POST("", in.DiscountController.Create)
		discountRouter.PUT("/:id", in.DiscountController.Update)
		discountRouter.GET("/list", in.DiscountController.GetDiscounts)
		discountRouter.GET("/:id", in.DiscountController.GetDiscountByID)
		discountRouter.GET("", in.DiscountController.GetDiscountByCode)
		discountRouter.PUT("/:id/status", in.DiscountController.SetStatus)
	}
}

func registerGenreRouter(root *gin.RouterGroup, in IRouter) {
	genreRouter := root.Group("/genre")
	{
		genreRouter.GET("", in.GenreController.GetGenres)
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

func registerStatisticRouter(root *gin.RouterGroup, in IRouter) {
	statisticRouter := root.Group("/statistic")
	{
		statisticRouter.GET("/movie-revenue", in.StatisticController.GetMovieRevenue)
		statisticRouter.GET("/cinema-revenue", in.StatisticController.GetCinemaRevenue)
		statisticRouter.GET("/combo", in.StatisticController.GetComboStatistics)
	}
}
