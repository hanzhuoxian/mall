package userserver

import (
	"fmt"

	pkgoptions "github.com/hanzhuoxian/mall/internal/pkg/options"
	"github.com/hanzhuoxian/mall/internal/pkg/server"
	"github.com/hanzhuoxian/mall/internal/userserver/config"
	"github.com/hanzhuoxian/mall/pkg/shutdown"
	"google.golang.org/grpc"
)

func provideMySQLOptions(cfg *config.Config) *pkgoptions.MySQLOptions {
	return cfg.MySQLOptions
}

func provideRedisOptions(cfg *config.Config) *pkgoptions.RedisOptions {
	return cfg.RedisOptions
}

func provideGracefulShutdown() *shutdown.GracefulShutdown {
	gs := shutdown.New()
	gs.AddShutdownManager(shutdown.NewPosixSignalManager())
	return gs
}

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

func provideGRPCServer(cfg *config.Config) *server.GRPCServer {
	addr := fmt.Sprintf("%s:%d", cfg.GRPCOptions.BindAddress, cfg.GRPCOptions.BindPort)
	var opts []grpc.ServerOption
	if cfg.GRPCOptions.MaxMsgSize > 0 {
		opts = append(opts,
			grpc.MaxRecvMsgSize(cfg.GRPCOptions.MaxMsgSize),
			grpc.MaxSendMsgSize(cfg.GRPCOptions.MaxMsgSize),
		)
	}
	return server.NewGRPCServer(addr, opts...)
}
