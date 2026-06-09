package app

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	flagHelp          = "help"
	flagHelpShorthand = "H"
)

// addHelpCommandFlag 向子命令的 FlagSet 添加 -H/--help flag。
func addHelpCommandFlag(usage string, fs *pflag.FlagSet) {
	fs.BoolP(
		flagHelp,
		flagHelpShorthand,
		false,
		fmt.Sprintf("Help for the %s command.", strings.Split(usage, " ")[0]),
	)
}

// helpCommand 返回一个自定义的 help 子命令，替换 cobra 默认的 help command，
// 支持 `<app> help [command]` 查看任意子命令的帮助信息。
func helpCommand(name string) *cobra.Command {
	return &cobra.Command{
		Use:   "help [command]",
		Short: "Help about any command",
		Long:  fmt.Sprintf("Help provides help for any command in the application.\nSimply type %s help [path to command] for full details.", name),
		RunE: func(c *cobra.Command, args []string) error {
			cmd, _, err := c.Root().Find(args)
			if cmd == nil || err != nil {
				c.Printf("Unknown help topic %#q\n", args)
				return c.Root().Usage()
			}
			cmd.InitDefaultHelpFlag()
			return cmd.Help()
		},
	}
}
