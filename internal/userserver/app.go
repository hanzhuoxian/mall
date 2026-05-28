package userserver

import (
	"github.com/hanzhuoxian/mall/internal/userserver/config"
	"github.com/hanzhuoxian/mall/internal/userserver/options"
	"github.com/hanzhuoxian/mall/pkg/app"
	"github.com/hanzhuoxian/mall/pkg/log"
)

const commandDesc = `The mall user server validates and configures data
for the api objects The Server services REST operations to do the api objects management.
`

func NewApp(basename string) *app.App {
	opts := options.NewOptions()
	a := app.NewApp(
		"Mall User Server",
		basename,
		app.WithOptions(opts),
		app.WithDescription(commandDesc),
		app.WithDefaultValidArgs(),
		app.WithRunFunc(run(opts)),
		app.WithNoConfig(),
	)
	return a
}

func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		log.Init(opts.LogOptions)
		defer log.Flush()

		cfg := config.CreateConfigFromOptions(opts)
		return Run(cfg)
	}
}
