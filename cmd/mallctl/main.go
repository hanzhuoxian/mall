package main

import (
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/hanzhuoxian/mall/internal/mallctl"
	"github.com/hanzhuoxian/mall/pkg/app"
	"github.com/hanzhuoxian/mall/pkg/logger"
)

func main() {
	if err := logger.Init(nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer func() { _ = zap.L().Sync() }()
	app.PrintWorkingDir()
	mallctl.NewApp("mallctl").Run()
}
