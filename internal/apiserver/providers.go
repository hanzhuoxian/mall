package apiserver

import (
	"github.com/hanzhuoxian/mall/internal/apiserver/config"
	"github.com/hanzhuoxian/mall/internal/apiserver/grpcclient"
	"github.com/hanzhuoxian/mall/internal/pkg/server"
	"github.com/hanzhuoxian/mall/pkg/shutdown"
)

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

// provideUserClient 从配置中提取 user service 地址并创建 gRPC 客户端。
func provideUserClient(cfg *config.Config) (*grpcclient.UserClient, error) {
	return grpcclient.NewUserClient(cfg.UserServiceOptions.Address())
}
