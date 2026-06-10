//go:build wireinject

package apiserver

import (
	"github.com/google/wire"
	"github.com/hanzhuoxian/mall/internal/apiserver/config"
	"github.com/hanzhuoxian/mall/internal/apiserver/controller"
)

// initApiServer is the Wire injector. wire gen will generate wire_gen.go from this.
func initApiServer(cfg *config.Config) (*apiserver, error) {
	wire.Build(
		// 基础设施
		provideGracefulShutdown,
		controller.ProviderSet,
		// server
		provideAPIServer,
		newApiServer,

		// grpc
		provideUserClient,

		// auth
		NewAutoAuth,
		NewJWTAuth,
	)
	return nil, nil
}
