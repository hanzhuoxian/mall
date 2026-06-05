package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	configFlagName  = "config"
	configFlagShort = "c"
	configFlagUsage = "Read configuration from specified `FILE`, support JSON, TOML, YAML, HCL, or Java properties formats."
)

// configFile 存储通过 --config/-c flag 指定的配置文件路径。
var configFile string

// init 将 --config/-c flag 注册到全局 pflag.CommandLine，使其在 AddConfigFlag 前即可被解析。
func init() {
	pflag.StringVarP(&configFile, configFlagName, configFlagShort, configFlagName, configFlagUsage)
}

// AddConfigFlag adds the config flag to the command.
func AddConfigFlag(basename string, fs *pflag.FlagSet) {
	fs.AddFlag(pflag.Lookup(configFlagName))

	viper.AutomaticEnv()
	viper.SetEnvPrefix(strings.Replace(strings.ToUpper(basename), "-", "_", -1))
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	cobra.OnInitialize(func() {
		if configFile != "" {
			viper.SetConfigFile(configFile)
		} else {
			viper.AddConfigPath(".")

			if names := strings.Split(basename, "-"); len(names) > 1 {
				viper.AddConfigPath(filepath.Join("/etc", names[0]))
			}

			viper.SetConfigName(basename)
		}

		if err := viper.ReadInConfig(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: failed to read configuration file(%s): %v\n", configFile, err)
			os.Exit(1)
		}
	})
}
