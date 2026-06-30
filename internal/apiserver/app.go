package apiserver

import (
	"context"
	"time"

	"github.com/hanzhuoxian/mall/internal/apiserver/options"
	"github.com/hanzhuoxian/mall/pkg/app"
	"github.com/hanzhuoxian/mall/pkg/logger"
	"github.com/hanzhuoxian/mall/pkg/telemetry"
)

// commandDesc 是用户服务命令行的长描述，展示在 --help 输出中。
const commandDesc = `The mall user server validates and configures data
for the api objects The Server services REST operations to do the api objects management.
`

// NewApp 创建并返回 bff 服务的 CLI 应用实例，配置好命令行参数和运行函数。
func NewApp(basename string) *app.App {
	options := options.NewOptions()
	return app.NewApp(
		"Mall Api Server",
		basename,
		app.WithOptions(options),
		app.WithDescription(commandDesc),
		app.WithRunFunc(run(options)),
		app.WithNoConfig(),
	)
}

func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		// 先初始化遥测，使全局 LoggerProvider 就绪，logger 才能桥接日志至 OTLP。
		// 服务名以 basename 兜底，可由 OTEL_SERVICE_NAME 环境变量覆盖。
		tp, err := telemetry.Init(context.Background(), opts.TelemetryOptions.Config(basename))
		if err != nil {
			return err
		}
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = tp.Shutdown(ctx)
		}()

		logger.Init(opts.LogOptions)
		defer logger.Flush()

		return Run(opts)
	}
}
