// Package middleware 提供 Gin 框架通用的 HTTP 中间件实现，
// 包括日志、请求上下文和注册的默认中间件集合。
package middleware

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-isatty"

	"github.com/hanzhuoxian/mall/pkg/log"
)

// defaultLogFormatter 默认的日志格式化函数，供 Logger 中间件使用。
// 它根据 gin 提供的参数生成结构化的一行日志文本。
var defaultLogFormatter = func(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency -= param.Latency % time.Second
	}

	return fmt.Sprintf("%s%3d%s - [%s] \"%v %s%s%s %s\" %s",
		// param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.ClientIP,
		param.Latency,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}

// Logger 返回使用默认配置的日志中间件。
// 可直接作为 Gin 路由组或全局中间件使用。
func Logger() gin.HandlerFunc {
	return LoggerWithConfig(GetLoggerConfig(nil, nil, nil))
}

// LoggerWithConfig 返回使用自定义配置的日志中间件，支持自定义格式化函数、输出目标和跳过路径。
// 参数 conf: 传入的 gin.LoggerConfig，若字段为空则使用合理的默认值。
func LoggerWithConfig(conf gin.LoggerConfig) gin.HandlerFunc {
	formatter := conf.Formatter
	if formatter == nil {
		formatter = defaultLogFormatter
	}

	out := conf.Output
	if out == nil {
		out = gin.DefaultWriter
	}

	notlogged := conf.SkipPaths

	isTerm := true

	if w, ok := out.(*os.File); !ok || os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd())) {
		isTerm = false
	}

	if isTerm {
		gin.ForceConsoleColor()
	}

	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			param := gin.LogFormatterParams{
				Request: c.Request,
				Keys:    c.Keys,
			}

			// Stop timer
			param.TimeStamp = time.Now()
			param.Latency = param.TimeStamp.Sub(start)

			param.ClientIP = c.ClientIP()
			param.Method = c.Request.Method
			param.StatusCode = c.Writer.Status()
			param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

			param.BodySize = c.Writer.Size()

			if raw != "" {
				path = path + "?" + raw
			}

			param.Path = path

			log.L(c).Info(formatter(param))
		}
	}
}

// GetLoggerConfig 构建并返回 gin.LoggerConfig，nil 参数将使用 gin 默认值。
// 常用于在程序初始化时组合自定义 formatter/output/skipPaths。
func GetLoggerConfig(formatter gin.LogFormatter, output io.Writer, skipPaths []string) gin.LoggerConfig {
	return gin.LoggerConfig{
		Formatter: formatter,
		Output:    output,
		SkipPaths: skipPaths,
	}
}
