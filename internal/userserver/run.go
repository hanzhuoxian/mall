package userserver

import (
	"fmt"

	"github.com/hanzhuoxian/mall/internal/userserver/config"
)

func Run(cfg *config.Config) error {
	server, err := createServer(cfg)
	if err != nil {
		return fmt.Errorf("failed to create server: %v", err)
	}
	return server.Run()
}
