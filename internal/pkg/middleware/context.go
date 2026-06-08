package middleware

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"

	"github.com/hanzhuoxian/mall/pkg/log"
)

// Context 返回一个中间件，将 requestID 和 username 注入到 gin.Context 中，
// 便于后续处理器通过 log.L(ctx) 自动携带这些字段输出结构化日志。
func Context() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(log.KeyRequestID, requestid.Get(c))
		c.Set(log.KeyUsername, c.GetString(log.KeyUsername))
		c.Next()
	}
}
