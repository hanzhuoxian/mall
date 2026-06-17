// Package middleware 提供 Gin 框架通用的 HTTP 中间件实现，
// 这里包含将 request id 与用户名注入到请求上下文的中间件。
package middleware

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/hanzhuoxian/mall/pkg/logger"
)

// Context 返回一个中间件，将 requestID 和 username 注入到 gin.Context 中，
// 便于后续处理器通过 logger.L(ctx) 自动携带这些字段输出结构化日志。
func Context() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(logger.KeyRequestID, requestid.Get(c))
		c.Set(logger.KeyUsername, c.GetString(logger.KeyUsername))
		c.Next()
	}
}
