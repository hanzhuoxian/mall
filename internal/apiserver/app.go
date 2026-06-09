package apiserver

import (
	"github.com/hanzhuoxian/mall/internal/apiserver/options"
	"github.com/hanzhuoxian/mall/pkg/app"
	"github.com/hanzhuoxian/mall/pkg/log"
)

// commandDesc 是用户服务命令行的长描述，展示在 --help 输出中。
const commandDesc = `The mall user server validates and configures data
for the api objects The Server services REST operations to do the api objects management.
`

// NewApp 创建并返回 bff 服务的 CLI 应用实例，配置好命令行参数和运行函数。
func NewApp(basename string) *app.App {
	options := options.NewOptions()
	return app.NewApp(
		basename,
		"Mall Api Server",
		app.WithDescription(commandDesc),
		app.WithRunFunc(run(options)),
		app.WithNoConfig(),
	)
}

func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		log.Init(opts.LogOptions)
		defer log.Flush()

		return Run(opts)
	}
}
