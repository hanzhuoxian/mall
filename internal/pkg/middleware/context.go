package middleware

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/hanzhuoxian/mall/pkg/log"
)

func Context() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(log.KeyRequestID, requestid.Get(c))
		c.Set(log.KeyUsername, c.GetString(log.KeyUsername))
		c.Next()
	}
}
