package controller

import (
	"TTCS/src/common/fault"
	"TTCS/src/common/log"
	"TTCS/src/common/ws"
	"TTCS/src/core/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

type WebSocketController struct {
	*BaseController
	authService *service.AuthService
	upgrader    websocket.Upgrader
	hub         *ws.Hub
}

func NewWebSocketController(baseController *BaseController, hub *ws.Hub, authService *service.AuthService) *WebSocketController {
	return &WebSocketController{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		BaseController: baseController,
		authService:    authService,
		hub:            hub,
	}
}

func (w *WebSocketController) HandleWebSocket(c *gin.Context) {
	caller := "WebSocketController.HandleWebSocket"
	ctx := c.Request.Context()

	room := c.Query("r")
	token := c.Query("t")

	if token == "" {
		err := fmt.Errorf("[%v] token is empty", caller)
		log.Error(ctx, err.Error())
		w.ServeErrResponse(c, fault.Wrapf(err, "[%v] token is empty", caller))
		return
	}
	user, err := w.authService.VerifyToken(ctx, token[7:])
	if err != nil {
		log.Error(ctx, "[%v] failed to verify token :%+v", caller, err)
		w.ServeErrResponse(c, err)
		return
	}

	conn, err := w.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Error(ctx, "Error upgrading to WebSocket: %+v", err)
		return
	}
	client := &ws.Client{
		ID:    user.ID,
		Conn:  conn,
		Route: room,
		Send:  make(chan []byte, 256),
	}

	// Đăng ký client vào Hub
	w.hub.Register <- client

	// Xử lý đọc tin nhắn từ client
	go client.Write()

	// Xử lý nhận tin nhắn từ client
	client.Read()

	w.hub.Register <- client
}
