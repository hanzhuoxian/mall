package app

import (
	"fmt"
	"os"

	cliflag "github.com/hanzhuoxian/mall/pkg/flag"
	"github.com/spf13/cobra"
)

type App struct {
	basename    string
	name        string
	description string
	options     CliOptions
	RunFunc     RunFunc
	silence     bool
	noVersion   bool
	noConfig    bool
	commands    []*Command
	args        cobra.PositionalArgs
	cmd         *cobra.Command
}

// Option defines optional parameters for initializing the application
// structure.
type Option func(*App)

// WithOptions defines the application's options.
func WithOptions(opt CliOptions) Option {
	return func(a *App) {
		a.options = opt
	}
}

// RunFunc defines the application's startup callback function.
type RunFunc func(basename string) error

// WithRunFunc defines the application's startup callback function.
func WithRunFunc(runFunc RunFunc) Option {
	return func(a *App) {
		a.RunFunc = runFunc
	}
}

// WithDescription defines the application's description.
func WithDescription(description string) Option {
	return func(a *App) {
		a.description = description
	}
}

// WithSilence defines whether to silence the application's output.
func WithSilence() Option {
	return func(a *App) {
		a.silence = true
	}
}

// WithNoVersion defines whether to disable the application's version command.
func WithNoVersion() Option {
	return func(a *App) {
		a.noVersion = true
	}
}

// WithNoConfig defines whether to disable the application's config command.
func WithNoConfig() Option {
	return func(a *App) {
		a.noConfig = true
	}
}

// WithValidArgs defines the application's valid args.
func WithValidArgs(args cobra.PositionalArgs) Option {
	return func(a *App) {
		a.args = args
	}
}

// WithDefaultValidArgs set default validation function to valid non-flag arguments.
func WithDefaultValidArgs() Option {
	return func(a *App) {
		a.args = func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}

			return nil
		}
	}
}

// NewApp creates a new application instance based on the given application name,
// binary name, and other options.
func NewApp(name string, basename string, opts ...Option) *App {
	a := &App{
		name:     name,
		basename: basename,
	}

	for _, o := range opts {
		o(a)
	}

	a.buildCommand()

	return a
}

func (a *App) buildCommand() {
	basename := FormatBaseName(a.basename)
	cmd := cobra.Command{
		Use:           basename,
		Short:         a.name,
		Long:          a.description,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          a.args,
	}

	cmd.Flags().SortFlags = true

	if len(a.commands) > 0 {
		for _, command := range a.commands {
			cmd.AddCommand(command.cobraCommand())
		}
		cmd.SetHelpCommand(helpCommand(basename))
	}

	if a.RunFunc != nil {
		cmd.RunE = a.runCommand
	}

	var namedFlagSet cliflag.NamedFlagSets
	if a.options != nil {
		namedFlagSet = a.options.Flags()
		fs := cmd.Flags()
		for _, f := range namedFlagSet.FlagSets {
			fs.AddFlagSet(f)
		}
	}

	a.cmd = &cmd
}

// Command returns cobra command instance inside the application.
func (a *App) Command() *cobra.Command {
	return a.cmd
}

// Run is used to launch the application.
func (a *App) Run() {
	if err := a.cmd.Execute(); err != nil {
		fmt.Printf("%v %v\n", "Error:", err)
		os.Exit(1)
	}
}

func (a *App) runCommand(cmd *cobra.Command, args []string) error {
	return nil
}
