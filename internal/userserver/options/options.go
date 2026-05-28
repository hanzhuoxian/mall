package options

import (
	pkgoptions "github.com/hanzhuoxian/mall/internal/pkg/options"
	"github.com/hanzhuoxian/mall/internal/pkg/server"
	"github.com/hanzhuoxian/mall/pkg/nflag"
)

type Options struct {
	ServerRunOptions       *pkgoptions.ServerRunOptions
	MySQLOptions           *pkgoptions.MySQLOptions
	GRPCOptions            *pkgoptions.GRPCOptions
	InsecureServingOptions *pkgoptions.InsecureServingOptions
}

func NewOptions() *Options {
	return &Options{
		ServerRunOptions:       pkgoptions.NewServerOptions(),
		MySQLOptions:           pkgoptions.NewMySQLOptions(),
		GRPCOptions:            pkgoptions.NewGRPCOptions(),
		InsecureServingOptions: pkgoptions.NewInsecureServingOptions(),
	}
}

func (o *Options) ApplyTo(c *server.Config) error {
	return nil
}

func (o *Options) Flags() (nfs nflag.NamedFlagSets) {
	o.ServerRunOptions.AddFlags(nfs.FlagSet("server"))
	o.GRPCOptions.AddFlags(nfs.FlagSet("grpc"))
	o.MySQLOptions.AddFlags(*nfs.FlagSet("mysql"))
	return nfs
}

func (o *Options) Validate() []error {
	var errs []error
	errs = append(errs, o.ServerRunOptions.Validate()...)
	errs = append(errs, o.GRPCOptions.Validate()...)
	errs = append(errs, o.MySQLOptions.Validate()...)
	errs = append(errs, o.InsecureServingOptions.Validate()...)
	return errs
}
