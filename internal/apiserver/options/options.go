// Package options 定义了 api 服务的所有命令行及配置文件选项。
package options

import (
	"fmt"

	pkgoptions "github.com/hanzhuoxian/mall/internal/pkg/options"
	"github.com/hanzhuoxian/mall/pkg/logger"
	"github.com/hanzhuoxian/mall/pkg/nflag"
)

// Options 汇总了 api 服务所有子系统的启动选项。
type Options struct {
	MySQLOptions           *pkgoptions.MySQLOptions           `json:"mysql"           mapstructure:"mysql"`
	InsecureServingOptions *pkgoptions.InsecureServingOptions `json:"insecure"        mapstructure:"insecure"`
	ServerRunOptions       *pkgoptions.ServerRunOptions       `json:"server"          mapstructure:"server"`
	UserServiceOptions     *pkgoptions.GRPCOptions            `json:"user-service"    mapstructure:"user-service"`
	RedisOptions           *pkgoptions.RedisOptions           `json:"redis"           mapstructure:"redis"`
	LogOptions             *logger.Options                    `json:"log"             mapstructure:"log"`
	TelemetryOptions       *pkgoptions.TelemetryOptions       `json:"otel"            mapstructure:"otel"`
	SandboxPayment         bool                               `json:"sandbox-payment" mapstructure:"sandbox-payment"`
}

// NewOptions 返回带有各子系统默认值的 Options 实例。
func NewOptions() *Options {
	insecureServing := pkgoptions.NewInsecureServingOptions()
	insecureServing.BindPort = 9091
	mysql := pkgoptions.NewMySQLOptions()
	mysql.Database = "mall"
	return &Options{
		MySQLOptions:           mysql,
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
	o.ServerRunOptions.AddFlags(nfs.FlagSet("server"))
	o.MySQLOptions.AddFlags(nfs.FlagSet("mysql"))
	nfs.FlagSet("commerce").
		BoolVar(&o.SandboxPayment, "sandbox-payment", false, "Enable simulated payments for local development only")
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
	if len(o.ServerRunOptions.JWTSecret) < 32 {
		errs = append(errs, fmt.Errorf("JWT secret must contain at least 32 bytes"))
	}
	if o.SandboxPayment && o.ServerRunOptions.Mode == "release" {
		errs = append(errs, fmt.Errorf("sandbox payments must not be enabled in release mode"))
	}
	errs = append(errs, o.ServerRunOptions.Validate()...)
	errs = append(errs, o.MySQLOptions.Validate()...)
	errs = append(errs, o.InsecureServingOptions.Validate()...)
	errs = append(errs, o.UserServiceOptions.Validate()...)
	errs = append(errs, o.RedisOptions.Validate()...)
	errs = append(errs, o.LogOptions.Validate()...)
	errs = append(errs, o.TelemetryOptions.Validate()...)
	return errs
}
