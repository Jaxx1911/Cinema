package middleware

import (
	"TTCS/src/common"
	"TTCS/src/common/log"
	"TTCS/src/present/httpui/controller"
	"fmt"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	*controller.BaseController
}

func NewAuthMiddleware(baseController *controller.BaseController) *AuthMiddleware {
	return &AuthMiddleware{
		baseController,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		caller := "AuthMiddleware.RequireAuth"

		token := c.Request.Header.Get("Authorization")
		if token == "" {
			err := fmt.Errorf("[%v] token is empty", caller)
			log.Error(ctx, err.Error())
			m.Error(c, common.ErrUnauthorized(c))
		}
	}
}
