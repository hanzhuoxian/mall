// Package middleware 提供 Gin 框架通用的 HTTP 中间件实现，
// 这里包含将 request id 与用户名注入到请求上下文的中间件。
package middleware

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/hanzhuoxian/mall/pkg/logger"
)

// Context 将请求元信息放入标准 context，供 HTTP、服务层和数据库日志共享。
func Context() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := requestid.Get(c)
		c.Set(logger.KeyRequestID, id)
		fields := []zap.Field{zap.String(logger.KeyRequestID, id)}
		if username := c.GetString(logger.KeyUsername); username != "" {
			fields = append(fields, zap.String(logger.KeyUsername, username))
		}
		ctx := logger.WithFields(c.Request.Context(), fields...)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
