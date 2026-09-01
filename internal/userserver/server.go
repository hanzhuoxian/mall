// Package userserver 实现了用户服务的核心逻辑，包括服务初始化、路由注册和优雅关闭。
package userserver

import (
	"errors"

	"github.com/hanzhuoxian/mall/internal/pkg/server"
	"github.com/hanzhuoxian/mall/internal/userserver/cache"
	"github.com/hanzhuoxian/mall/internal/userserver/service"
	"github.com/hanzhuoxian/mall/internal/userserver/store"
	"github.com/hanzhuoxian/mall/pkg/shutdown"
)

// userServer 是用户服务的核心结构，持有所有基础设施实例（HTTP 服务、gRPC 服务、数据层等）。
type userServer struct {
	gs           *shutdown.GracefulShutdown
	grpcServer   *server.GRPCServer
	apiServer    *server.APIServer
	storeFactory store.Factory
	cacheFactory cache.Factory
	svc          service.Service
}

// preparedUserServer 是完成路由注册和关闭回调注册后可直接运行的服务实例。
type preparedUserServer struct {
	*userServer
}

// newUserServer 创建并返回一个新的 userServer 实例，由 Wire 注入所有依赖。
func newUserServer(
	gs *shutdown.GracefulShutdown,
	apiServer *server.APIServer,
	grpcServer *server.GRPCServer,
	sf store.Factory,
	cf cache.Factory,
	svc service.Service,
) *userServer {
	return &userServer{
		gs:           gs,
		apiServer:    apiServer,
		grpcServer:   grpcServer,
		storeFactory: sf,
		cacheFactory: cf,
		svc:          svc,
	}
}

// PrepareRun 完成路由注册和关闭回调注册，返回可直接调用 Run 的 preparedUserServer。
func (u *userServer) PrepareRun() *preparedUserServer {
	installRoutes(u.grpcServer, u.svc)

	// 关闭顺序不可调整：必须先停止接收新请求并排空在途请求（HTTP、gRPC），
	// 再关闭它们依赖的 MySQL / Redis。反过来会让在途请求拿到已关闭的连接池，
	// 优雅关停就失去了意义。
	//
	// GracefulShutdown 并发执行各个 callback，因此有先后依赖的清理必须写在同一个 callback 内。
	u.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(s string) error {
		u.apiServer.Close()
		grpcErr := u.grpcServer.Close()

		storeErr := u.storeFactory.Close()
		cacheErr := u.cacheFactory.Close()

		return errors.Join(grpcErr, storeErr, cacheErr)
	}))
	return &preparedUserServer{u}
}

// Run 启动优雅关闭管理器、gRPC 服务和 HTTP API 服务，并阻塞直到关闭完成。
func (p preparedUserServer) Run() error {
	if err := p.gs.Start(); err != nil {
		return err
	}

	if err := p.grpcServer.Run(); err != nil {
		return err
	}

	if err := p.apiServer.Run(); err != nil {
		return err
	}
	<-p.gs.Done()
	return nil
}
