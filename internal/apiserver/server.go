package apiserver

import (
	"github.com/hanzhuoxian/mall/internal/apiserver/controller"
	"github.com/hanzhuoxian/mall/internal/pkg/middleware"
	"github.com/hanzhuoxian/mall/internal/pkg/middleware/auth"
	"github.com/hanzhuoxian/mall/internal/pkg/server"
	"github.com/hanzhuoxian/mall/pkg/shutdown"
)

type apiserver struct {
	gs           *shutdown.GracefulShutdown
	apiServer    *server.APIServer
	controllers  *controller.Controllers
	authStrategy middleware.AuthStrategy
	jwtStrategy  auth.JWTStrategy
}

type preparedApiServer struct {
	*apiserver
}

func newApiServer(gs *shutdown.GracefulShutdown, s *server.APIServer, c *controller.Controllers, authStrategy middleware.AuthStrategy, jwtStrategy auth.JWTStrategy) (*apiserver, error) {
	return &apiserver{
		gs:           gs,
		apiServer:    s,
		controllers:  c,
		authStrategy: authStrategy,
		jwtStrategy:  jwtStrategy,
	}, nil
}

func (a *apiserver) PrepareRun() *preparedApiServer {
	installRoutes(a.apiServer.Engine, a.controllers, a.authStrategy, a.jwtStrategy)

	a.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(s string) error {
		a.apiServer.Close()
		return nil
	}))
	return &preparedApiServer{a}
}

func (pa *preparedApiServer) Run() error {
	if err := pa.gs.Start(); err != nil {
		return err
	}
	if err := pa.apiServer.Run(); err != nil {
		return err
	}
	<-pa.gs.Done()
	return nil
}
