package apiserver

import (
	"fmt"

	"github.com/hanzhuoxian/mall/internal/apiserver/config"
	"github.com/hanzhuoxian/mall/internal/apiserver/options"
)

// Run 从 Options 构建配置，初始化并启动 api server。
func Run(opts *options.Options) error {
	cfg := config.CreateConfigFromOptions(opts)
	srv, err := initApiServer(cfg)
	if err != nil {
		return fmt.Errorf("failed to init server: %w", err)
	}
	return srv.PrepareRun().Run()
}
