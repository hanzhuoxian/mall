// Package nflag 在 pflag 基础上提供了命名分组 FlagSet（NamedFlagSets）能力，
// 支持按分组顺序打印 flag 帮助信息并自适应终端宽度。
package nflag

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/pflag"
)

// GlobalFlagSetName 是全局通用 flag 分组的保留名称，version 和 config flag 默认注册到该分组。
const GlobalFlagSetName = "global"

// NamedFlagSets 管理多个具名 pflag.FlagSet，Order 维护插入顺序以保证打印时顺序一致。
type NamedFlagSets struct {
	FlagSets map[string]*pflag.FlagSet
	Order    []string
}

// FlagSet returns a FlagSet with the given name
func (nfs *NamedFlagSets) FlagSet(name string) *pflag.FlagSet {
	if nfs.FlagSets == nil {
		nfs.FlagSets = make(map[string]*pflag.FlagSet)
	}
	if nfs.FlagSets[name] == nil {
		nfs.Order = append(nfs.Order, name)
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
			_, _ = fmt.Fprint(w, strings.Join(lines[:len(lines)-1], "\n"))
			_, _ = fmt.Fprintln(w)
		} else {
			_, _ = fmt.Fprint(w, buf.String())
		}
	}
}

// WordSepNormalizeFunc changes all flags that contain "_" separators.
func WordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		return pflag.NormalizedName(strings.ReplaceAll(name, "_", "-"))
	}
	return pflag.NormalizedName(name)
}

// InitFlags normalizes, parses, then logs the command line flags.
func InitFlags(flags *pflag.FlagSet) {
	flags.SetNormalizeFunc(WordSepNormalizeFunc)
	flags.AddGoFlagSet(flag.CommandLine)
}
