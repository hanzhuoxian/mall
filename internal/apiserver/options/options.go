// Package options 定义了 api 服务的所有命令行及配置文件选项。
package options

import (
	pkgoptions "github.com/hanzhuoxian/mall/internal/pkg/options"
	"github.com/hanzhuoxian/mall/pkg/logger"
	"github.com/hanzhuoxian/mall/pkg/nflag"
)

// Options 汇总了 api 服务所有子系统的启动选项。
type Options struct {
	InsecureServingOptions *pkgoptions.InsecureServingOptions `json:"insecure" mapstructure:"insecure"`
	ServerRunOptions       *pkgoptions.ServerRunOptions       `json:"server" mapstructure:"server"`
	UserServiceOptions     *pkgoptions.GRPCOptions            `json:"user-service" mapstructure:"user-service"`
	RedisOptions           *pkgoptions.RedisOptions           `json:"redis" mapstructure:"redis"`
	LogOptions             *logger.Options                    `json:"log"      mapstructure:"log"`
	TelemetryOptions       *pkgoptions.TelemetryOptions       `json:"otel" mapstructure:"otel"`
}

// NewOptions 返回带有各子系统默认值的 Options 实例。
func NewOptions() *Options {
	insecureServing := pkgoptions.NewInsecureServingOptions()
	insecureServing.BindPort = 9090
	return &Options{
		InsecureServingOptions: insecureServing,
		ServerRunOptions:       pkgoptions.NewServerOptions(),
		UserServiceOptions:     pkgoptions.NewGRPCOptions(),
		RedisOptions:           pkgoptions.NewRedisOptions(),
		LogOptions:             logger.NewOptions(),
		TelemetryOptions:       pkgoptions.NewTelemetryOptions(),
	}
}

// Flags 返回按子系统分组的命名 FlagSet。
func (o *Options) Flags() (nfs nflag.NamedFlagSets) {
	o.InsecureServingOptions.AddFlags(nfs.FlagSet("insecure"))
	o.UserServiceOptions.AddFlags(nfs.FlagSet("user-service"))
	o.RedisOptions.AddFlags(nfs.FlagSet("redis"))
	o.LogOptions.AddFlags(nfs.FlagSet("log"))
	o.TelemetryOptions.AddFlags(nfs.FlagSet("otel"))
	return nfs
}

// Validate 依次校验各子系统选项合法性。
func (o *Options) Validate() []error {
	var errs []error
	errs = append(errs, o.InsecureServingOptions.Validate()...)
	errs = append(errs, o.UserServiceOptions.Validate()...)
	errs = append(errs, o.RedisOptions.Validate()...)
	errs = append(errs, o.LogOptions.Validate()...)
	errs = append(errs, o.TelemetryOptions.Validate()...)
	return errs
}
