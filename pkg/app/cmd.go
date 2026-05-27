package app

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// Command is a sub command structure of a cli application.
// It is recommended that a command be created with the app.NewCommand()
// function.
type Command struct {
	usage    string
	desc     string
	options  CliOptions
	commands []*Command
	runFunc  RunCommandFunc
}

// CommonOptions are options that are common to all commands.
type CommonOptions func(*Command)

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
func NewCommand(name string, usage string, desc string, options ...CommonOptions) *Command {
	cmd := &Command{
		usage: usage,
		desc:  desc,
	}

	for _, option := range options {
		option(cmd)
	}

	return cmd
}

func (c *Command) cobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   c.usage,
		Short: c.desc,
	}
	cmd.SetOutput(os.Stdout)
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
		for _, f := range c.options.Flags().FlagSets {
			cmd.Flags().AddFlagSet(f)
		}
		// c.options.AddFlags(cmd.Flags())
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

func (c *Command) RunCommand(args []string) {
	if c.runFunc == nil {
		return
	}

	if err := c.runFunc(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// AddCommand adds sub command to the application.
func (a *App) AddCommand(cmd *Command) {
	a.commands = append(a.commands, cmd)
}

// AddCommands adds multiple sub commands to the application.
func (a *App) AddCommands(cmds ...*Command) {
	a.commands = append(a.commands, cmds...)
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
