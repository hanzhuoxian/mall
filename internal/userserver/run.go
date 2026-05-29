package userserver

import (
	"fmt"

	"github.com/hanzhuoxian/mall/internal/userserver/config"
)

func Run(cfg *config.Config) error {
	srv, err := initUserServer(cfg)
	if err != nil {
		return fmt.Errorf("failed to init server: %w", err)
	}
	return srv.PrepareRun().Run()
}
