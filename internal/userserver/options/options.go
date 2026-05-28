package options

import (
	pkgoptions "github.com/hanzhuoxian/mall/internal/pkg/options"
	"github.com/hanzhuoxian/mall/internal/pkg/server"
	"github.com/hanzhuoxian/mall/pkg/nflag"
)

type Options struct {
	ServerRunOptions *pkgoptions.ServerRunOptions
}

func NewOptions() *Options {
	return &Options{
		ServerRunOptions: pkgoptions.NewServerOptions(),
	}
}

func (o *Options) ApplyTo(c *server.Config) error {
	return nil
}

func (o *Options) Flags() (nfs nflag.NamedFlagSets) {
	return nfs
}
