package http

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func GenRsUID() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Request-Id", uuid.NewV4().String())
		c.Next()
	}
}
