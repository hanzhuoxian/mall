package mallctl

import (
	"github.com/hanzhuoxian/mall/pkg/app"
)

// commandDesc 是 mallctl 命令行工具的长描述，展示在 --help 输出中。
const commandDesc = `mallctl is the command-line management tool for the mall platform.
It provides administrative operations for managing mall API resources.
`

// NewApp 创建并返回 mallctl 命令行工具的应用实例。
func NewApp(basename string) *app.App {
	a := app.NewApp(
		basename,
		"Mall Control Tool",
		app.WithDescription(commandDesc),
		app.WithNoConfig(),
		app.WithCommands(newMigrateCmd()),
	)

	return a
}
