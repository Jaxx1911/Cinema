package middleware

import (
	"TTCS/src/common/fault"
	"TTCS/src/common/log"
	"TTCS/src/core/service"
	"TTCS/src/present/httpui/controller"
	"fmt"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	*controller.BaseController
	authService *service.AuthService
}

func NewAuthMiddleware(baseController *controller.BaseController, authService *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		baseController,
		authService,
	}
}

func (a *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		caller := "AuthMiddleware.RequireAuth"

		token := c.Request.Header.Get("Authorization")
		if token == "" {
			err := fmt.Errorf("[%v] token is empty", caller)
			log.Error(ctx, err.Error())
			a.ServeErrResponse(c, fault.Wrapf(err, "[%v] token is empty", caller))
		}
		user, err := a.authService.VerifyToken(ctx, token[7:])
		if err != nil {
			log.Error(ctx, "[%v] failed to verify token :%+v", caller, err)
			a.ServeErrResponse(c, err)
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
