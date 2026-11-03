// Package logger 提供项目的 zap 初始化、配置及请求上下文关联。
package logger

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"go.uber.org/zap/zapcore"
)

// 命令行 flag 名称常量，供 AddFlags 注册使用。
const (
	flagLevel             = "log.level"
	flagDisableCaller     = "log.disable-caller"
	flagDisableStacktrace = "log.disable-stacktrace"
	flagFormat            = "log.format"
	flagEnableColor       = "log.enable-color"
	flagOutputPaths       = "log.output-paths"
	flagErrorOutputPaths  = "log.error-output-paths"
	flagDevelopment       = "log.development"
	flagName              = "log.name"

	consoleFormat = "console"
	jsonFormat    = "json"
)

// Options 包含日志相关的所有配置项，支持 JSON 反序列化和 mapstructure 映射。
type Options struct {
	Level             string   `json:"level"              mapstructure:"level"`
	Format            string   `json:"format"             mapstructure:"format"`
	Name              string   `json:"name"               mapstructure:"name"`
	OutputPaths       []string `json:"output-paths"       mapstructure:"output-paths"`
	ErrorOutputPaths  []string `json:"error-output-paths" mapstructure:"error-output-paths"`
	DisableCaller     bool     `json:"disable-caller"     mapstructure:"disable-caller"`
	DisableStacktrace bool     `json:"disable-stacktrace" mapstructure:"disable-stacktrace"`
	EnableColor       bool     `json:"enable-color"       mapstructure:"enable-color"`
	Development       bool     `json:"development"        mapstructure:"development"`
}

// NewOptions 返回带有合理默认值的 Options 实例：Info 级别、console 格式、输出到 stdout/stderr。
func NewOptions() (o *Options) {
	return &Options{
		Level:             zapcore.InfoLevel.String(),
		DisableCaller:     false,
		DisableStacktrace: false,
		Format:            consoleFormat,
		EnableColor:       false,
		Development:       false,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
	}
}

// Validate 校验配置项合法性，返回所有错误；New 在创建 logger 前执行此校验。
func (o *Options) Validate() []error {
	var errors []error
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
		errors = append(errors, err)
	}

	format := strings.ToLower(o.Format)
	if format != consoleFormat && format != jsonFormat {
		errors = append(errors, fmt.Errorf("not a valid log format %q", format))
	}
	return errors
}

// AddFlags adds flags for log to the specified FlagSet object.
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Level, flagLevel, o.Level, "Minimum log output `LEVEL`.")
	fs.BoolVar(&o.DisableCaller, flagDisableCaller, o.DisableCaller, "Disable output of caller information in the log.")
	fs.BoolVar(&o.DisableStacktrace, flagDisableStacktrace,
		o.DisableStacktrace, "Disable the log to record a stack trace for all messages at or above panic level.")
	fs.StringVar(&o.Format, flagFormat, o.Format, "Log output `FORMAT`, support console or json format.")
	fs.BoolVar(&o.EnableColor, flagEnableColor, o.EnableColor, "Enable output ansi colors in console format logs.")
	fs.StringSliceVar(&o.OutputPaths, flagOutputPaths, o.OutputPaths, "Output paths of log.")
	fs.StringSliceVar(&o.ErrorOutputPaths, flagErrorOutputPaths, o.ErrorOutputPaths, "Error output paths of log.")
	fs.BoolVar(
		&o.Development,
		flagDevelopment,
		o.Development,
		"Development puts the logger in development mode, which makes DPanic logs panic.",
	)
	fs.StringVar(&o.Name, flagName, o.Name, "The name of the logger.")
}

// String 将配置序列化为 JSON 字符串，便于调试输出。
func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}
