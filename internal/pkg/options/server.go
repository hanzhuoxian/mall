// Package options 定义了各基础设施（HTTP、gRPC、MySQL、Redis 等）的命令行/配置选项结构。
package options

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"

	"github.com/hanzhuoxian/mall/internal/pkg/server"
)

// ServerRunOptions 包含 HTTP API 服务器的基本运行选项，如模式、健康检查和中间件列表。
type ServerRunOptions struct {
	Mode            string        `json:"mode"             mapstructure:"mode"`
	Middlewares     []string      `json:"middlewares"      mapstructure:"middlewares"`
	Healthz         bool          `json:"healthz"          mapstructure:"healthz"`
	JWTSecret       string        `json:"jwt-secret"       mapstructure:"jwt-secret"`
	ShutdownTimeout time.Duration `json:"shutdown-timeout" mapstructure:"shutdown-timeout"`
}

// NewServerOptions 返回与 server.NewConfig() 默认值一致的 ServerRunOptions 实例。
func NewServerOptions() *ServerRunOptions {
	defaults := server.NewConfig()
	return &ServerRunOptions{
		Mode:            defaults.Mode,
		Healthz:         defaults.Healthz,
		Middlewares:     defaults.Middlewares,
		ShutdownTimeout: defaults.ShutdownTimeout,
	}
}

// ApplyTo 将选项值写入 server.Config，供服务器初始化时使用。
func (s *ServerRunOptions) ApplyTo(c *server.Config) error {
	c.Mode = s.Mode
	c.Healthz = s.Healthz
	c.Middlewares = s.Middlewares
	c.ShutdownTimeout = s.ShutdownTimeout
	return nil
}

// Validate 校验选项合法性。
func (s *ServerRunOptions) Validate() []error {
	var errors []error
	if s.ShutdownTimeout <= 0 {
		errors = append(errors, fmt.Errorf("--server.shutdown-timeout must be positive, got %s", s.ShutdownTimeout))
	}
	return errors
}

// AddFlags 向指定 FlagSet 注册服务器运行参数的命令行 flag。
func (s *ServerRunOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.Mode, "server.mode", s.Mode, "server mode, Supported server mode: debug, test, release")
	fs.BoolVar(&s.Healthz, "server.healthz", s.Healthz, "enable healthz")
	fs.StringSliceVar(&s.Middlewares, "server.middlewares", s.Middlewares, "server middlewares")
	fs.StringVar(&s.JWTSecret, "server.jwt-secret", s.JWTSecret, "JWT signing secret")
	fs.DurationVar(&s.ShutdownTimeout, "server.shutdown-timeout", s.ShutdownTimeout,
		"maximum duration to wait for in-flight requests to finish during graceful shutdown")
}
