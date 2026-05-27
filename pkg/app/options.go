package app

import "github.com/hanzhuoxian/mall/pkg/nflag"

// CliOptions defines flags for a particular type of component.
type CliOptions interface {
	// AddFlags adds flags to the specified FlagSet object.
	// AddFlags(fs *pflag.FlagSet)
	Flags() (fss nflag.NamedFlagSets)
	Validate() []error
}

// ConfigurableOptions abstracts configuration options for reading parameters
// from a configuration file.
type ConfigurableOptions interface {
	// ApplyFlags parsing parameters from the command line or configuration file
	// to the options instance.
	ApplyFlags() []error
}

// CompleteableOptions abstracts options which can be completed.
type CompleteableOptions interface {
	Complete() error
}

// PrintableOptions abstracts options which can be printed.
type PrintableOptions interface {
	String() string
}
