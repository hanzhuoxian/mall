package options

import (
	pkgoptions "github.com/hanzhuoxian/mall/internal/pkg/options"
	"github.com/hanzhuoxian/mall/internal/pkg/server"
	"github.com/hanzhuoxian/mall/pkg/log"
	"github.com/hanzhuoxian/mall/pkg/nflag"
)

type Options struct {
	ServerRunOptions       *pkgoptions.ServerRunOptions
	MySQLOptions           *pkgoptions.MySQLOptions
	GRPCOptions            *pkgoptions.GRPCOptions
	InsecureServingOptions *pkgoptions.InsecureServingOptions
	LogOptions             *log.Options
	RedisOptions           *pkgoptions.RedisOptions
}

func NewOptions() *Options {
	return &Options{
		ServerRunOptions:       pkgoptions.NewServerOptions(),
		MySQLOptions:           pkgoptions.NewMySQLOptions(),
		GRPCOptions:            pkgoptions.NewGRPCOptions(),
		InsecureServingOptions: pkgoptions.NewInsecureServingOptions(),
		RedisOptions:           pkgoptions.NewRedisOptions(),
	}
}

func (o *Options) ApplyTo(c *server.Config) error {
	return nil
}

func (o *Options) Flags() (nfs nflag.NamedFlagSets) {
	o.ServerRunOptions.AddFlags(nfs.FlagSet("server"))
	o.GRPCOptions.AddFlags(nfs.FlagSet("grpc"))
	o.MySQLOptions.AddFlags(nfs.FlagSet("mysql"))
	o.InsecureServingOptions.AddFlags(nfs.FlagSet("insecure serving"))
	o.RedisOptions.AddFlags(nfs.FlagSet("redis"))
	return nfs
}

func (o *Options) Validate() []error {
	var errs []error
	errs = append(errs, o.ServerRunOptions.Validate()...)
	errs = append(errs, o.GRPCOptions.Validate()...)
	errs = append(errs, o.MySQLOptions.Validate()...)
	errs = append(errs, o.InsecureServingOptions.Validate()...)
	errs = append(errs, o.RedisOptions.Validate()...)
	return errs
}
