package options

import pkgoptions "github.com/hanzhuoxian/mall/internal/pkg/options"

type Options struct {
	ServerRunOptions *pkgoptions.ServerRunOptions
}

func NewOptions() *Options {
	return &Options{
		ServerRunOptions: pkgoptions.NewServerOptions(),
	}
}
