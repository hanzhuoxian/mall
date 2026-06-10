//go:build wireinject

package userserver

import (
	"github.com/google/wire"
	"github.com/hanzhuoxian/mall/internal/userserver/cache/redis"
	"github.com/hanzhuoxian/mall/internal/userserver/config"
	"github.com/hanzhuoxian/mall/internal/userserver/service"
	"github.com/hanzhuoxian/mall/internal/userserver/store/mysql"
)

// initUserServer is the Wire injector. wire gen will generate wire_gen.go from this.
func initUserServer(cfg *config.Config) (*userServer, error) {
	wire.Build(
		// config 拆分
		provideMySQLOptions,
		provideRedisOptions,
		// 基础设施
		provideGracefulShutdown,
		provideAPIServer,
		provideGRPCServer,
		// store / cache
		mysql.ProviderSet,
		redis.ProviderSet,
		// service
		service.ProviderSet,
		// server
		newUserServer,
	)
	return nil, nil
}
