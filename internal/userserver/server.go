package userserver

import (
	"fmt"

	"github.com/hanzhuoxian/mall/internal/pkg/options"
	"github.com/hanzhuoxian/mall/internal/pkg/server"
	"github.com/hanzhuoxian/mall/internal/userserver/cache"
	"github.com/hanzhuoxian/mall/internal/userserver/cache/redis"
	"github.com/hanzhuoxian/mall/internal/userserver/config"
	"github.com/hanzhuoxian/mall/internal/userserver/store"
	"github.com/hanzhuoxian/mall/internal/userserver/store/mysql"
	"github.com/hanzhuoxian/mall/pkg/shutdown"
	"google.golang.org/grpc"
)

type userServer struct {
	gs         *shutdown.GracefulShutdown
	grpcServer *server.GRPCServer
	apiServer  *server.APIServer
}

type preparedUserServer struct {
	*userServer
}

type ExtraConfig struct {
	Addr           string
	MaxMessageSize int
	mysqlOptions   options.MySQLOptions
	redisOptions   options.RedisOptions
}

func createServer(cfg *config.Config) (*userServer, error) {
	gs := shutdown.New()

	gs.AddShutdownManager(shutdown.NewPosixSignalManager())

	c, err := buildConfig(cfg)
	if err != nil {
		return nil, err
	}
	apiServer, err := c.Complete().New()
	if err != nil {
		return nil, err
	}

	ec, err := buildExtraConfig(cfg)
	if err != nil {
		return nil, err
	}

	grpcServer, err := ec.Complete().New()
	if err != nil {
		return nil, err
	}

	return &userServer{
		gs:         gs,
		apiServer:  apiServer,
		grpcServer: grpcServer,
	}, nil
}

func (u *userServer) PrepareRun() *preparedUserServer {
	initRouter(u.apiServer.Engine)
	u.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(s string) error {
		if s, _ := mysql.GetMySQLFactoryOr(nil); s != nil {
			s.Close()
		}
		if c, _ := redis.GetCacheFactoryOr(nil); c != nil {
			c.Close()
		}
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

func buildConfig(cfg *config.Config) (c *server.Config, lastErr error) {
	c = server.NewConfig()
	if lastErr = cfg.ServerRunOptions.ApplyTo(c); lastErr != nil {
		return
	}
	if lastErr = cfg.InsecureServingOptions.ApplyTo(c); lastErr != nil {
		return
	}
	return
}

func buildExtraConfig(cfg *config.Config) (*ExtraConfig, error) {
	return &ExtraConfig{
		Addr:           fmt.Sprintf("%s:%d", cfg.GRPCOptions.BindAddress, cfg.GRPCOptions.BindPort),
		MaxMessageSize: cfg.GRPCOptions.MaxMsgSize,
		mysqlOptions:   *cfg.MySQLOptions,
		redisOptions:   *cfg.RedisOptions,
	}, nil
}

type completedExtraConfig struct {
	*ExtraConfig
}

func (e *ExtraConfig) Complete() *completedExtraConfig {
	if e.Addr == "" {
		e.Addr = "127.0.0.1:8081"
	}

	return &completedExtraConfig{e}
}

func (ce *completedExtraConfig) New() (*server.GRPCServer, error) {
	opts := []grpc.ServerOption{}
	if ce.MaxMessageSize > 0 {
		opts = append(opts,
			grpc.MaxRecvMsgSize(ce.MaxMessageSize),
			grpc.MaxSendMsgSize(ce.MaxMessageSize),
		)
	}
	s, err := mysql.GetMySQLFactoryOr(&ce.mysqlOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to init mysql: %w", err)
	}
	store.Set(s)

	c, err := redis.GetCacheFactoryOr(&ce.redisOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to init redis: %w", err)
	}
	cache.Set(c)

	return server.NewGRPCServer(ce.Addr, opts...), nil
}
