package options

import (
	"errors"
	"os"

	"github.com/spf13/pflag"

	"github.com/hanzhuoxian/mall/pkg/telemetry"
	"github.com/hanzhuoxian/mall/pkg/version"
)

// errInvalidSamplerRatio 表示采样比例超出 [0,1] 合法范围。
var errInvalidSamplerRatio = errors.New("otel.sampler-ratio must be in range [0,1]")

// TelemetryOptions 包含 OpenTelemetry（Trace/Metric/Log）的命令行/配置选项，
// 通过 OTLP/gRPC 协议将遥测数据导出至 collector。Endpoint 留空时回退到
// 标准的 OTEL_EXPORTER_OTLP_* 环境变量，从而支持纯环境变量配置。
type TelemetryOptions struct {
	Endpoint     string  `json:"endpoint,omitempty"      mapstructure:"endpoint"`
	Environment  string  `json:"environment,omitempty"   mapstructure:"environment"`
	SamplerRatio float64 `json:"sampler-ratio,omitempty" mapstructure:"sampler-ratio"`
	Enabled      bool    `json:"enabled"                 mapstructure:"enabled"`
	Insecure     bool    `json:"insecure"                mapstructure:"insecure"`
}

// NewTelemetryOptions 返回带有合理默认值的 TelemetryOptions 实例。
// 默认关闭，开发环境下连接本地 collector（localhost:4317，明文，全采样）。
func NewTelemetryOptions() *TelemetryOptions {
	return &TelemetryOptions{
		Enabled:      true,
		Endpoint:     "localhost:4317",
		Environment:  "dev",
		SamplerRatio: 1.0,
		Insecure:     true,
	}
}

// Validate 校验遥测选项合法性，主要检查采样比例范围。
func (o *TelemetryOptions) Validate() []error {
	var errs []error
	if o.Enabled && (o.SamplerRatio < 0 || o.SamplerRatio > 1) {
		errs = append(errs, errInvalidSamplerRatio)
	}
	return errs
}

// AddFlags 向指定 FlagSet 注册 OpenTelemetry 相关的命令行 flag。
func (o *TelemetryOptions) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&o.Enabled, "otel.enabled", o.Enabled,
		"Enable OpenTelemetry traces, metrics and logs export via OTLP/gRPC.")

	fs.StringVar(&o.Endpoint, "otel.endpoint", o.Endpoint,
		"OTLP/gRPC collector endpoint (host:port). If blank, OTEL_EXPORTER_OTLP_* env vars are used.")

	fs.StringVar(&o.Environment, "otel.environment", o.Environment,
		"Deployment environment name reported as the deployment.environment.name resource attribute.")

	fs.Float64Var(&o.SamplerRatio, "otel.sampler-ratio", o.SamplerRatio,
		"Trace sampling ratio in [0,1]. 1 samples every trace.")

	fs.BoolVar(&o.Insecure, "otel.insecure", o.Insecure,
		"Use a plaintext (non-TLS) gRPC connection to the OTLP collector.")
}

// Config 将选项转换为 telemetry.Config。fallbackName 为服务名兜底（通常是进程 basename），
// 可被 OpenTelemetry 标准的 OTEL_SERVICE_NAME 环境变量覆盖——服务名不再单独提供命令行 flag。
func (o *TelemetryOptions) Config(fallbackName string) telemetry.Config {
	serviceName := fallbackName
	if name := os.Getenv("OTEL_SERVICE_NAME"); name != "" {
		serviceName = name
	}
	return telemetry.Config{
		Enabled:        o.Enabled,
		ServiceName:    serviceName,
		ServiceVersion: version.Get().GitVersion,
		Environment:    o.Environment,
		Endpoint:       o.Endpoint,
		SamplerRatio:   o.SamplerRatio,
		Insecure:       o.Insecure,
	}
}
