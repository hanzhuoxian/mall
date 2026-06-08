// Package config 定义了用户服务的运行时配置结构，并提供从选项构建配置的工厂函数。
package config

import "github.com/hanzhuoxian/mall/internal/userserver/options"

// Config 是用户服务的运行时配置，内嵌 Options 以复用所有命令行/配置文件选项。
type Config struct {
	*options.Options
}

// CreateConfigFromOptions 将解析后的 Options 包装为 Config 实例。
func CreateConfigFromOptions(opts *options.Options) *Config {
	return &Config{opts}
}
