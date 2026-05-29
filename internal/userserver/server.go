package userserver

import (
	"github.com/hanzhuoxian/mall/internal/pkg/server"
	"github.com/hanzhuoxian/mall/internal/userserver/cache"
	"github.com/hanzhuoxian/mall/internal/userserver/controller"
	"github.com/hanzhuoxian/mall/internal/userserver/store"
	"github.com/hanzhuoxian/mall/pkg/shutdown"
)

type userServer struct {
	gs           *shutdown.GracefulShutdown
	grpcServer   *server.GRPCServer
	apiServer    *server.APIServer
	storeFactory store.Factory
	cacheFactory cache.Factory
	controllers  *controller.Controllers
}

type preparedUserServer struct {
	*userServer
}

func newUserServer(
	gs *shutdown.GracefulShutdown,
	apiServer *server.APIServer,
	grpcServer *server.GRPCServer,
	sf store.Factory,
	cf cache.Factory,
	ctrls *controller.Controllers,
) *userServer {
	return &userServer{
		gs:           gs,
		apiServer:    apiServer,
		grpcServer:   grpcServer,
		storeFactory: sf,
		cacheFactory: cf,
		controllers:  ctrls,
	}
}

func (u *userServer) PrepareRun() *preparedUserServer {
	initRouter(u.apiServer.Engine, u.controllers)
	u.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(s string) error {
		u.storeFactory.Close()
		u.cacheFactory.Close()
		u.apiServer.Close()
		u.grpcServer.Close()
		return nil
	}))
	return &preparedUserServer{u}
}

func (p preparedUserServer) Run() error {
	if err := p.gs.Start(); err != nil {
		return err
	}

	go p.grpcServer.Run()

	if err := p.apiServer.Run(); err != nil {
		return err
	}
	<-p.gs.Done()
	return nil
}
