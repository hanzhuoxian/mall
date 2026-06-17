package userserver

import (
	"github.com/hanzhuoxian/mall/internal/userserver/config"
	"github.com/hanzhuoxian/mall/internal/userserver/options"
	"github.com/hanzhuoxian/mall/pkg/app"
	"github.com/hanzhuoxian/mall/pkg/logger"
)

// commandDesc 是用户服务命令行的长描述，展示在 --help 输出中。
const commandDesc = `The mall user server validates and configures data
for the api objects The Server services REST operations to do the api objects management.
`

// NewApp 创建并返回用户服务的 CLI 应用实例，配置好命令行参数和运行函数。
func NewApp(basename string) *app.App {
	opts := options.NewOptions()
	return app.NewApp(
		"Mall User Server",
		basename,
		app.WithOptions(opts),
		app.WithDescription(commandDesc),
		app.WithDefaultValidArgs(),
		app.WithRunFunc(run(opts)),
		app.WithNoConfig(),
	)
}

// run 返回一个 RunFunc，在应用启动时初始化日志、构建配置并运行服务。
func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		logger.Init(opts.LogOptions)
		defer logger.Flush()

		cfg := config.CreateConfigFromOptions(opts)
		return Run(cfg)
	}
}
