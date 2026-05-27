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

	return &server{}, nil
}
func (s *server) Run() error {
	fmt.Println("hello user server!")
	return nil
}
