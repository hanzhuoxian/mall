package apiserver

import (
	"context"
	"time"

	"github.com/hanzhuoxian/mall/internal/apiserver/controller"
	marketservice "github.com/hanzhuoxian/mall/internal/marketplace/service"
	"github.com/hanzhuoxian/mall/internal/pkg/middleware"
	"github.com/hanzhuoxian/mall/internal/pkg/middleware/auth"
	"github.com/hanzhuoxian/mall/internal/pkg/server"
	"github.com/hanzhuoxian/mall/pkg/logger"
	"github.com/hanzhuoxian/mall/pkg/shutdown"
)

type apiserver struct {
	marketplace  *marketservice.Service
	gs           *shutdown.GracefulShutdown
	apiServer    *server.APIServer
	controllers  *controller.Controllers
	authStrategy middleware.AuthStrategy
	jwtStrategy  auth.JWTStrategy
}

type preparedApiServer struct {
	*apiserver
}

func newApiServer(
	gs *shutdown.GracefulShutdown,
	s *server.APIServer,
	c *controller.Controllers,
	authStrategy middleware.AuthStrategy,
	jwtStrategy auth.JWTStrategy,
	marketplace *marketservice.Service,
) (*apiserver, error) {
	return &apiserver{
		marketplace:  marketplace,
		gs:           gs,
		apiServer:    s,
		controllers:  c,
		authStrategy: authStrategy,
		jwtStrategy:  jwtStrategy,
	}, nil
}

func (a *apiserver) PrepareRun() *preparedApiServer {
	installRoutes(a.apiServer.Engine, a.controllers, a.authStrategy, a.jwtStrategy)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := a.marketplace.ExpireOrders(ctx); err != nil && ctx.Err() == nil {
					logger.Errorf("expire orders: %v", err)
				}
			}
		}
	}()
	a.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(s string) error {
		cancel()
		<-done
		a.apiServer.Close()
		return a.marketplace.Close()
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
