// Package server 提供 HTTP（Gin）和 gRPC 服务器的抽象封装，
// 包含配置结构、服务器生命周期管理（启动、路由初始化、优雅关闭）及健康检查等能力。
package server

import (
	"net"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	// RecommendedHomeDir 是 mall 在用户主目录下的默认配置目录名。
	RecommendedHomeDir = ".mall"

	// RecommendedEnvPrefix 是 mall 环境变量的统一前缀，用于 viper 自动读取。
	RecommendedEnvPrefix = "MALL"
)

// DefaultShutdownTimeout 是优雅关闭 HTTP 服务器时等待在途请求处理完成的默认上限。
const DefaultShutdownTimeout = 60 * time.Second

// Config 是 APIServer 的完整配置，包含运行模式、中间件列表、可选功能开关及 TLS/非 TLS 监听信息。
type Config struct {
	SecureServing   *SecureServingInfo
	InsecureServing *InsecureServingInfo
	Jwt             *JwtInfo
	Mode            string
	Middlewares     []string
	ShutdownTimeout time.Duration
	Healthz         bool
	EnableProfiling bool
	EnableMetrics   bool
}

// CertKey 持有 TLS 证书和私钥的文件路径。
type CertKey struct {
	CertFile string
	KeyFile  string
}

// SecureServingInfo 包含 HTTPS 服务的监听地址、端口及证书信息。
type SecureServingInfo struct {
	BindAddress string
	CertKey     CertKey
	BindPort    int
}

// Address 返回格式化的 HTTPS 监听地址字符串（host:port）。
func (s *SecureServingInfo) Address() string {
	return net.JoinHostPort(s.BindAddress, strconv.Itoa(s.BindPort))
}

// InsecureServingInfo 包含 HTTP 服务的监听地址（已合并 host:port）。
type InsecureServingInfo struct {
	Address string
}

// JwtInfo 包含 JWT 认证所需的域名、密钥及令牌有效期配置。
type JwtInfo struct {
	Realm      string
	Key        string
	Timeout    time.Duration
	MaxRefresh time.Duration
}

// NewConfig 返回带有合理默认值的 Config 实例（debug 模式、健康检查开启、pprof 和 metrics 开启）。
func NewConfig() *Config {
	return &Config{
		Healthz:         true,
		Mode:            gin.DebugMode,
		Middlewares:     []string{},
		ShutdownTimeout: DefaultShutdownTimeout,
		EnableProfiling: true,
		EnableMetrics:   true,
		Jwt: &JwtInfo{
			Realm:      "mall",
			Timeout:    1 * time.Hour,
			MaxRefresh: 1 * time.Hour,
		},
	}
}

// CompletedConfig 是经过完整性检查的配置，只能通过 Config.Complete() 获取。
type CompletedConfig struct {
	*Config
}

// Complete 对配置进行完整性填充，返回可用于创建服务器的 CompletedConfig。
func (cfg *Config) Complete() CompletedConfig {
	if cfg.ShutdownTimeout <= 0 {
		cfg.ShutdownTimeout = DefaultShutdownTimeout
	}
	return CompletedConfig{cfg}
}

// New 根据已完成配置创建并初始化 APIServer 实例（设置模式、注册中间件和内置 API）。
func (c CompletedConfig) New() (*APIServer, error) {
	gin.SetMode(c.Mode)

	s := &APIServer{
		SecureServingInfo:   c.SecureServing,
		InsecureServingInfo: c.InsecureServing,
		ShutdownTimeout:     c.ShutdownTimeout,
		healthz:             c.Healthz,
		enableMetrics:       c.EnableMetrics,
		enableProfiling:     c.EnableProfiling,
		middlewares:         c.Middlewares,
		Engine:              gin.New(),
	}

	initAPIServer(s)
	return s, nil
}
