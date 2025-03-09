package middleware

import (
	"TTCS/src/common/log"
	"github.com/gin-gonic/gin"
)

func Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		log.Info(c, "path: [%v], status: [%v], method:[%v]", c.Request.URL.Path, c.Writer.Status(), c.Request.Method)
	}
}
