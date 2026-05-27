package version

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/pflag"
)

type versionValue int

// 版本值的定义
const (
	VersionFalse versionValue = iota
	VersionTrue
	VersionRaw
)

const strRawVersion string = "raw"

// IsBoolFlag indicates that this flag can be used as a boolean flag (i.e., --version will be treated as --version=true).
func (v *versionValue) IsBoolFlag() bool {
	return true
}

// Get returns the current version value.
func (v *versionValue) Get() any {
	return v
}

// Set sets the version value based on the input string. It accepts "raw" to set VersionRaw, and parses boolean values for VersionTrue and VersionFalse.
func (v *versionValue) Set(s string) error {
	if s == strRawVersion {
		*v = VersionRaw
		return nil
	}
	boolValue, err := strconv.ParseBool(s)
	if boolValue {
		*v = VersionTrue
	} else {
		*v = VersionFalse
	}
	return err
}

// String returns a string representation of the version value.
func (v *versionValue) String() string {
	if *v == VersionRaw {
		return strRawVersion
	}
	return fmt.Sprintf("%v", bool(*v == VersionTrue))
}

// Type returns a string name for this Option type.  This will be used in help and error messages.
func (v *versionValue) Type() string {
	return "version"
}

// VersionVar defines a flag with the specified name and usage string.
func VersionVar(p *versionValue, name string, value versionValue, usage string) {
	*p = value
	pflag.Var(p, name, usage)
	// "--version" will be treated as "--version=true"
	pflag.Lookup(name)
}

// Version wraps the VersionVar function.
func Version(name string, value versionValue, usage string) *versionValue {
	p := new(versionValue)
	VersionVar(p, name, value, usage)
	return p
}

// VersionFlagName is the name of the version flag.
const VersionFlagName = "version"

var versionFlag = Version(VersionFlagName, VersionFalse, "Print version information and quit")

// AddFlags adds the version flag to the given FlagSet.
func AddFlags(fs *pflag.FlagSet) {
	fs.AddFlag(pflag.Lookup(VersionFlagName))
}

// PrintAndExitIfRequested will check if the -version flag was passed
// and, if so, print the version and exit.
func PrintAndExitIfRequested() {
	switch *versionFlag {
	case VersionRaw:
		fmt.Printf("%#v\n", Get())
		os.Exit(0)
	case VersionTrue:
		fmt.Printf("%s\n", Get())
		os.Exit(0)
	}
}
