package app

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/pflag"
)

// NamedFlagSet is a wrapper around pflag.FlagSet
type NamedFlagSets struct {
	Order    []string
	FlagSets map[string]*pflag.FlagSet
}

// FlagSet returns a FlagSet with the given name
func (nfs *NamedFlagSets) FlagSet(name string) *pflag.FlagSet {
	if nfs.FlagSets == nil {
		nfs.FlagSets = make(map[string]*pflag.FlagSet)
	}
	if nfs.FlagSets[name] == nil {
		nfs.FlagSets[name] = pflag.NewFlagSet(name, pflag.ExitOnError)
	}
	return nfs.FlagSets[name]
}

// PrintSections prints the flag sets in the given NamedFlagSets
func PrintSections(w io.Writer, fss NamedFlagSets, cols int) {
	for _, name := range fss.Order {
		fs := fss.FlagSets[name]
		if !fs.HasFlags() {
			continue
		}

		wideFS := pflag.NewFlagSet("", pflag.ExitOnError)
		wideFS.AddFlagSet(fs)

		var zzz string
		if cols > 24 {
			zzz = strings.Repeat("z", cols-24)
			wideFS.Int(zzz, 0, strings.Repeat("z", cols-24))
		}

		var buf bytes.Buffer
		fmt.Fprintf(&buf, "\n%s flags:\n\n%s", strings.ToUpper(name[:1])+name[1:], wideFS.FlagUsagesWrapped(cols))

		if cols > 24 {
			i := strings.Index(buf.String(), zzz)
			lines := strings.Split(buf.String()[:i], "\n")
			fmt.Fprint(w, strings.Join(lines[:len(lines)-1], "\n"))
			fmt.Fprintln(w)
		} else {
			fmt.Fprint(w, buf.String())
		}
	}
}
