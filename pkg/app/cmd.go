package app

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hanzhuoxian/mall/pkg/nflag"
	"github.com/hanzhuoxian/mall/pkg/term"
)

// Command is a sub command structure of a cli application.
// It is recommended that a command be created with the app.NewCommand()
// function.
type Command struct {
	usage    string
	desc     string
	options  CliOptions
	runFunc  RunCommandFunc
	commands []*Command
}

// CommonOptions are options that are common to all commands.
type CommonOptions func(*Command)

// RunCommandFunc defines the application's command startup callback function.
type RunCommandFunc func(args []string) error

// WithCommandOptions adds options to a command.
func WithCommandOptions(options CliOptions) CommonOptions {
	return func(cmd *Command) {
		cmd.options = options
	}
}

// WithCommandRunFunc adds a run function to a command.
func WithCommandRunFunc(runFunc RunCommandFunc) CommonOptions {
	return func(cmd *Command) {
		cmd.runFunc = runFunc
	}
}

// NewCommand creates a new command.
func NewCommand(usage string, desc string, options ...CommonOptions) *Command {
	cmd := &Command{
		usage: usage,
		desc:  desc,
	}

	for _, option := range options {
		option(cmd)
	}

	return cmd
}

// cobraCommand 将 Command 转换为 cobra.Command，注册子命令、flag 和 help flag。
func (c *Command) cobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   c.usage,
		Short: c.desc,
	}
	cmd.SetOut(os.Stdout)
	cmd.Flags().SortFlags = false
	if len(c.commands) > 0 {
		for _, command := range c.commands {
			cmd.AddCommand(command.cobraCommand())
		}
	}

	if c.runFunc != nil {
		cmd.Run = c.runCommand
	}

	if c.options != nil {
		namedFlagSets := c.options.Flags()
		for _, f := range namedFlagSets.FlagSets {
			cmd.Flags().AddFlagSet(f)
		}
		cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
			cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\n\nUsage:\n  %s\n\n", cmd.Short, cmd.UseLine())
			nflag.PrintSections(cmd.OutOrStdout(), namedFlagSets, cols)
		})
	}
	addHelpCommandFlag(c.usage, cmd.Flags())

	return cmd
}

// AddCommand adds a sub command to a command.
func (c *Command) AddCommand(command *Command) {
	c.commands = append(c.commands, command)
}

// AddCommands adds multiple sub commands to a command.
func (c *Command) AddCommands(commands ...*Command) {
	c.commands = append(c.commands, commands...)
}

// runCommand 是子命令的 cobra Run 实现，执行失败时打印错误并以非零状态退出。
func (c *Command) runCommand(cmd *cobra.Command, args []string) {
	if c.runFunc == nil {
		return
	}

	if err := c.runFunc(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// FormatBaseName is formatted as an executable file name under different
// operating systems according to the given name.
func FormatBaseName(basename string) string {
	// Make case-insensitive and strip executable suffix if present
	if runtime.GOOS == "windows" {
		basename = strings.ToLower(basename)
		basename = strings.TrimSuffix(basename, ".exe")
	}

	return basename
}
