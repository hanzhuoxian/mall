// Package middleware 提供 Gin 框架使用的通用 HTTP 中间件集合。
package middleware

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

// Middlewares 是按名称索引的全局中间件注册表，服务器可按名称按需加载。
var Middlewares = defaultMiddlewares()

// defaultMiddlewares 返回内置的默认中间件集合，包括 recovery、requestid 和 logger。
func defaultMiddlewares() map[string]gin.HandlerFunc {
	return map[string]gin.HandlerFunc{
		"recovery":  gin.Recovery(),
		"requestid": requestid.New(),
		"logger":    Logger(),
	}
}
