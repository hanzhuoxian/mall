// Package options 定义了用户服务的所有命令行及配置文件选项，
// 汇聚各子系统（HTTP、gRPC、MySQL、Redis、日志）的配置结构。
package options

import (
	pkgoptions "github.com/hanzhuoxian/mall/internal/pkg/options"
	"github.com/hanzhuoxian/mall/internal/pkg/server"
	"github.com/hanzhuoxian/mall/pkg/logger"
	"github.com/hanzhuoxian/mall/pkg/nflag"
)

// Options 汇总了用户服务所有子系统的启动选项。
type Options struct {
	ServerRunOptions       *pkgoptions.ServerRunOptions       `json:"server"   mapstructure:"server"`
	MySQLOptions           *pkgoptions.MySQLOptions           `json:"mysql"    mapstructure:"mysql"`
	GRPCOptions            *pkgoptions.GRPCOptions            `json:"grpc"     mapstructure:"grpc"`
	InsecureServingOptions *pkgoptions.InsecureServingOptions `json:"insecure" mapstructure:"insecure"`
	LogOptions             *logger.Options                    `json:"log"      mapstructure:"log"`
	RedisOptions           *pkgoptions.RedisOptions           `json:"redis"    mapstructure:"redis"`
	TelemetryOptions       *pkgoptions.TelemetryOptions       `json:"otel"     mapstructure:"otel"`
}

// NewOptions 返回带有各子系统默认值的 Options 实例。
func NewOptions() *Options {
	mysqlOpt := pkgoptions.NewMySQLOptions()
	mysqlOpt.Database = "mall_user"
	return &Options{
		ServerRunOptions:       pkgoptions.NewServerOptions(),
		MySQLOptions:           mysqlOpt,
		GRPCOptions:            pkgoptions.NewGRPCOptions(),
		InsecureServingOptions: pkgoptions.NewInsecureServingOptions(),
		RedisOptions:           pkgoptions.NewRedisOptions(),
		TelemetryOptions:       pkgoptions.NewTelemetryOptions(),
	}
}

// ApplyTo 将选项应用到服务器配置，当前为空占位。
func (o *Options) ApplyTo(c *server.Config) error {
	return nil
}

// Flags 返回按子系统分组的命名 FlagSet，供 CLI 框架注册命令行参数使用。
func (o *Options) Flags() (nfs nflag.NamedFlagSets) {
	o.ServerRunOptions.AddFlags(nfs.FlagSet("server"))
	o.GRPCOptions.AddFlags(nfs.FlagSet("grpc"))
	o.MySQLOptions.AddFlags(nfs.FlagSet("mysql"))
	o.InsecureServingOptions.AddFlags(nfs.FlagSet("insecure serving"))
	o.RedisOptions.AddFlags(nfs.FlagSet("redis"))
	o.TelemetryOptions.AddFlags(nfs.FlagSet("otel"))
	return nfs
}

// Validate 依次校验各子系统选项的合法性，返回所有错误列表。
func (o *Options) Validate() []error {
	var errs []error
	errs = append(errs, o.ServerRunOptions.Validate()...)
	errs = append(errs, o.GRPCOptions.Validate()...)
	errs = append(errs, o.MySQLOptions.Validate()...)
	errs = append(errs, o.InsecureServingOptions.Validate()...)
	errs = append(errs, o.RedisOptions.Validate()...)
	errs = append(errs, o.TelemetryOptions.Validate()...)
	return errs
}
