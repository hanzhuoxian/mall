package config

import "github.com/hanzhuoxian/mall/internal/userserver/options"

type Config struct {
	*options.Options
}

func CreateConfigFromOptions(opts *options.Options) *Config {
	return &Config{opts}
}
