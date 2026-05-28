package userserver

import (
	"fmt"

	pkgserver "github.com/hanzhuoxian/mall/internal/pkg/server"
	"github.com/hanzhuoxian/mall/internal/userserver/config"
)

type server struct {
	grpcServer *pkgserver.GRPCServer
	apiServer  *pkgserver.APIServer
}

func createServer(cfg *config.Config) (*server, error) {
	c, err := buildConfig(cfg)
	if err != nil {
		return nil, err
	}
	c.Mode()
	return &server{}, nil
}
func (s *server) Run() error {
	fmt.Println("hello user server!")
	return nil
}

func buildConfig(cfg *config.Config) (c *pkgserver.Config, lastErr error) {
	c = pkgserver.NewConfig()
	if lastErr = cfg.ServerRunOptions.ApplyTo(c); lastErr != nil {
		return
	}
	return
}
