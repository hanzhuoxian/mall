package userserver

import (
	"fmt"

	"github.com/hanzhuoxian/mall/internal/userserver/config"
)

// Run 根据配置初始化用户服务并启动，是服务的主要运行入口。
func Run(cfg *config.Config) error {
	srv, err := initUserServer(cfg)
	if err != nil {
		return fmt.Errorf("failed to init server: %w", err)
	}
	return srv.PrepareRun().Run()
}
