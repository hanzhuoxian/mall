package userserver

import (
	"fmt"

	"github.com/hanzhuoxian/mall/internal/pkg/options"
	"github.com/hanzhuoxian/mall/internal/pkg/server"
	"github.com/hanzhuoxian/mall/internal/userserver/config"
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

	u.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(s string) error {
		u.apiServer.Close()
		u.grpcServer.Close()
		return nil
	}))
	return &preparedUserServer{u}
}

func (p preparedUserServer) Run() error {
	go p.grpcServer.Run()

	return p.apiServer.Run()
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
	return server.NewGRPCServer(ce.Addr, opts...), nil
}
