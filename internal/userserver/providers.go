// Package userserver 中的 providers.go 提供 Wire 所需的依赖构建函数，
// 负责从配置中提取各基础设施选项并构造对应实例。
package userserver

import (
	"fmt"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	"github.com/hanzhuoxian/mall/internal/pkg/middleware"
	pkgoptions "github.com/hanzhuoxian/mall/internal/pkg/options"
	"github.com/hanzhuoxian/mall/internal/pkg/server"
	"github.com/hanzhuoxian/mall/internal/userserver/config"
	"github.com/hanzhuoxian/mall/pkg/shutdown"
)

// provideMySQLOptions 从全局配置中提取 MySQL 选项，供 Wire 注入数据库连接使用。
func provideMySQLOptions(cfg *config.Config) *pkgoptions.MySQLOptions {
	return cfg.MySQLOptions
}

// provideRedisOptions 从全局配置中提取 Redis 选项，供 Wire 注入缓存连接使用。
func provideRedisOptions(cfg *config.Config) *pkgoptions.RedisOptions {
	return cfg.RedisOptions
}

// provideGracefulShutdown 创建优雅关闭管理器并注册 POSIX 信号监听器（SIGINT/SIGTERM）。
func provideGracefulShutdown() *shutdown.GracefulShutdown {
	gs := shutdown.New()
	gs.AddShutdownManager(shutdown.NewPosixSignalManager())
	return gs
}

// provideAPIServer 根据配置构建 HTTP API 服务器实例。
func provideAPIServer(cfg *config.Config) (*server.APIServer, error) {
	c := server.NewConfig()
	if err := cfg.ServerRunOptions.ApplyTo(c); err != nil {
		return nil, err
	}
	if err := cfg.InsecureServingOptions.ApplyTo(c); err != nil {
		return nil, err
	}
	return c.Complete().New()
}

// provideGRPCServer 根据配置构建 gRPC 服务器实例，若配置了消息大小限制则同时设置收发限制。
func provideGRPCServer(cfg *config.Config) *server.GRPCServer {
	addr := fmt.Sprintf("%s:%d", cfg.GRPCOptions.BindAddress, cfg.GRPCOptions.BindPort)
	opts := []grpc.ServerOption{
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			middleware.ValidationUnaryInterceptor,
		),
	}
	if cfg.GRPCOptions.MaxMsgSize > 0 {
		opts = append(opts,
			grpc.MaxRecvMsgSize(cfg.GRPCOptions.MaxMsgSize),
			grpc.MaxSendMsgSize(cfg.GRPCOptions.MaxMsgSize),
		)
	}
	return server.NewGRPCServer(addr, opts...)
}
