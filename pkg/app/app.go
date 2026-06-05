package app

import (
	"fmt"
	"os"

	"github.com/hanzhuoxian/mall/pkg/log"
	"github.com/hanzhuoxian/mall/pkg/nflag"
	"github.com/hanzhuoxian/mall/pkg/term"
	"github.com/hanzhuoxian/mall/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var progressMessage = "==>"

type App struct {
	basename    string
	name        string
	description string
	options     CliOptions
	runFunc     RunFunc
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

// RunFunc defines the application's startup callback function.
type RunFunc func(basename string) error

// WithOptions defines the application's options.
func WithOptions(opt CliOptions) Option {
	return func(a *App) {
		a.options = opt
	}
}

// WithRunFunc defines the application's startup callback function.
func WithRunFunc(runFunc RunFunc) Option {
	return func(a *App) {
		a.runFunc = runFunc
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
	app := &App{
		name:     name,
		basename: basename,
	}

	for _, o := range opts {
		o(app)
	}

	app.buildCommand()

	return app
}

func (app *App) buildCommand() {
	printWorkingDir()
	basename := FormatBaseName(app.basename)
	cmd := cobra.Command{
		Use:           basename,
		Short:         app.name,
		Long:          app.description,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          app.args,
	}

	cmd.Flags().SortFlags = true
	nflag.InitFlags(cmd.Flags())

	if len(app.commands) > 0 {
		for _, command := range app.commands {
			cmd.AddCommand(command.cobraCommand())
		}
		cmd.SetHelpCommand(helpCommand(basename))
	}

	if app.runFunc != nil {
		cmd.RunE = app.runCommand
	}

	var namedFlagSet nflag.NamedFlagSets
	if app.options != nil {
		namedFlagSet = app.options.Flags()
		fs := cmd.Flags()
		for _, f := range namedFlagSet.FlagSets {
			fs.AddFlagSet(f)
		}
	}
	if !app.noVersion {
		version.AddFlags(namedFlagSet.FlagSet(nflag.GlobalFlagSetName))
	}

	if !app.noConfig {
		AddConfigFlag(basename, namedFlagSet.FlagSet(nflag.GlobalFlagSetName))
	}

	nflag.AddGlobalFlags(namedFlagSet.FlagSet(nflag.GlobalFlagSetName), cmd.Name())
	cmd.Flags().AddFlagSet(namedFlagSet.FlagSet(nflag.GlobalFlagSetName))

	addCmdTemplate(&cmd, namedFlagSet)

	app.cmd = &cmd
}

func printWorkingDir() {
	wd, _ := os.Getwd()
	log.Infof("%v WorkingDir: %s", progressMessage, wd)
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
	if !a.noVersion {
		version.PrintAndExitIfRequested()
	}

	if !a.noConfig {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := viper.Unmarshal(&a.options); err != nil {
			return err
		}
	}
	if a.options != nil {
		if err := a.applyOptionRules(); err != nil {
			return err
		}
	}
	if a.runFunc != nil {
		return a.runFunc(a.basename)
	}

	return nil
}

func (a *App) applyOptionRules() error {
	if completeableOptions, ok := a.options.(CompleteableOptions); ok {
		if err := completeableOptions.Complete(); err != nil {
			return err
		}
	}

	if errs := a.options.Validate(); len(errs) != 0 {
		return fmt.Errorf("invalid options: %v", errs)
	}

	return nil
}

func addCmdTemplate(cmd *cobra.Command, namedFlagSets nflag.NamedFlagSets) {
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		cols, _, _ := term.TerminalSize(cmd.OutOrStderr())
		fmt.Fprintf(cmd.OutOrStderr(), "Usage:\n  %s\n\n", cmd.UseLine())
		nflag.PrintSections(cmd.OutOrStderr(), namedFlagSets, cols)
		return nil
	})

	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\nUsage:\n  %s\n\n", cmd.Long, cmd.UseLine())
		nflag.PrintSections(cmd.OutOrStdout(), namedFlagSets, cols)
	})
}
