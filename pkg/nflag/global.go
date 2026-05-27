package nflag

import (
	"flag"
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

// AddGlobalFlags adds global flags to the specified flagset.
func AddGlobalFlags(fs *pflag.FlagSet, name string) {
	fs.BoolP("help", "h", false, fmt.Sprintf("help for %s", name))
}

// normalize normalizes the flag name by replacing underscores with hyphens.
func normalize(s string) string {
	return strings.ReplaceAll(s, "_", "-")
}

// Register registers a flag from the global flagset to the specified flagset.
func Register(fs *pflag.FlagSet, globalName string) {
	if f := flag.CommandLine.Lookup(globalName); f != nil {
		pflagFlag := pflag.PFlagFromGoFlag(f)
		pflagFlag.Name = normalize(pflagFlag.Name)
		fs.AddFlag(pflagFlag)
	} else {
		panic(fmt.Sprintf("failed to find flag in global flagset (flag): %s", globalName))
	}
}
