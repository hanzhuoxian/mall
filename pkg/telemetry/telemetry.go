// Package telemetry 封装了 OpenTelemetry SDK 的初始化逻辑，统一管理 Trace、Metric、Log
// 三类信号的导出（均通过 OTLP/gRPC 协议）。它负责构建 resource、配置全局 Provider 和
// 文本传播器，并返回一个可在优雅关闭时调用的统一 Shutdown 入口。
//
// 各类自动埋点库（otelgin、otelgrpc、otelgorm、redisotel、otelzap）均依赖此处设置的
// 全局 Provider，因此应在应用启动早期、构建其他基础设施之前调用 Init。
package telemetry

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	otellog "go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

// Config 描述 OpenTelemetry 的运行时配置。Endpoint 为空时，导出器会回退到标准的
// OTEL_EXPORTER_OTLP_* 环境变量（默认 localhost:4317），从而支持纯环境变量配置。
type Config struct {
	// ServiceName 标识当前服务，写入 resource 的 service.name 属性。
	ServiceName string
	// ServiceVersion 为服务版本，写入 service.version 属性。
	ServiceVersion string
	// Environment 为部署环境（如 dev/prod），写入 deployment.environment.name 属性。
	Environment string
	// Endpoint 为 OTLP/gRPC collector 地址（host:port，不含 scheme）。留空则使用环境变量。
	Endpoint string
	// SamplerRatio 为基于 TraceID 的采样比例，取值 [0,1]。1 表示全采样。
	SamplerRatio float64
	// Insecure 为 true 时使用明文 gRPC 连接（开发环境常用）。
	Insecure bool
	// Enabled 为总开关，false 时 Init 直接返回空 Provider，不安装任何导出器。
	Enabled bool
}

// Provider 持有已初始化的各类 SDK Provider 及其关闭函数，用于在程序退出时统一释放资源。
type Provider struct {
	shutdownFuncs []func(context.Context) error
}

// Init 根据 cfg 初始化全局 TracerProvider、MeterProvider、LoggerProvider 及文本传播器。
// 若 cfg.Enabled 为 false，则返回一个 Shutdown 为空操作的 Provider，所有埋点退化为 no-op。
// 任意子系统初始化失败时，会回滚已创建的资源并返回错误。
func Init(ctx context.Context, cfg Config) (*Provider, error) {
	p := &Provider{}
	if !cfg.Enabled {
		return p, nil
	}

	// 文本传播器使用 W3C TraceContext + Baggage，保证跨服务（HTTP/gRPC）链路连续。
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	res, err := newResource(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := p.initTracer(ctx, cfg, res); err != nil {
		_ = p.Shutdown(ctx)
		return nil, fmt.Errorf("init tracer provider: %w", err)
	}

	if err := p.initMeter(ctx, cfg, res); err != nil {
		_ = p.Shutdown(ctx)
		return nil, fmt.Errorf("init meter provider: %w", err)
	}

	if err := p.initLogger(ctx, cfg, res); err != nil {
		_ = p.Shutdown(ctx)
		return nil, fmt.Errorf("init logger provider: %w", err)
	}

	return p, nil
}

// newResource 构建包含服务标识属性的 OTel resource，并自动合并环境检测到的属性。
func newResource(ctx context.Context, cfg Config) (*resource.Resource, error) {
	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithProcessPID(),
		resource.WithProcessExecutableName(),
		resource.WithProcessExecutablePath(),
		resource.WithProcessOwner(),
		resource.WithProcessRuntimeName(),
		resource.WithProcessRuntimeVersion(),
		resource.WithProcessRuntimeDescription(),
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			semconv.DeploymentEnvironmentName(cfg.Environment),
		),
	)
	if errors.Is(err, resource.ErrPartialResource) || errors.Is(err, resource.ErrSchemaURLConflict) {
		return res, nil
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

// initTracer 创建 OTLP/gRPC trace 导出器与批处理 TracerProvider，并注册为全局 Provider。
func (p *Provider) initTracer(ctx context.Context, cfg Config, res *resource.Resource) error {
	opts := []otlptracegrpc.Option{}
	if cfg.Endpoint != "" {
		opts = append(opts, otlptracegrpc.WithEndpoint(cfg.Endpoint))
	}
	if cfg.Insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}

	exporter, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(cfg.SamplerRatio))),
	)
	otel.SetTracerProvider(tp)
	p.shutdownFuncs = append(p.shutdownFuncs, tp.Shutdown)
	return nil
}

// initMeter 创建 OTLP/gRPC metric 导出器与周期性 MeterProvider，注册为全局 Provider，
// 并启动 Go runtime 指标采集（GC、goroutine、内存等）。
func (p *Provider) initMeter(ctx context.Context, cfg Config, res *resource.Resource) error {
	opts := []otlpmetricgrpc.Option{}
	if cfg.Endpoint != "" {
		opts = append(opts, otlpmetricgrpc.WithEndpoint(cfg.Endpoint))
	}
	if cfg.Insecure {
		opts = append(opts, otlpmetricgrpc.WithInsecure())
	}

	exporter, err := otlpmetricgrpc.New(ctx, opts...)
	if err != nil {
		return err
	}

	reader := sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(15*time.Second))
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(reader),
	)
	otel.SetMeterProvider(mp)
	p.shutdownFuncs = append(p.shutdownFuncs, mp.Shutdown)

	// 采集 Go 运行时指标，使用刚注册的全局 MeterProvider。
	if err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(15 * time.Second)); err != nil {
		return err
	}
	return nil
}

// initLogger 创建 OTLP/gRPC log 导出器与批处理 LoggerProvider，并注册为全局 LoggerProvider，
// 供 otelzap 桥接器将应用日志导出至 collector。
func (p *Provider) initLogger(ctx context.Context, cfg Config, res *resource.Resource) error {
	opts := []otlploggrpc.Option{}
	if cfg.Endpoint != "" {
		opts = append(opts, otlploggrpc.WithEndpoint(cfg.Endpoint))
	}
	if cfg.Insecure {
		opts = append(opts, otlploggrpc.WithInsecure())
	}

	exporter, err := otlploggrpc.New(ctx, opts...)
	if err != nil {
		return err
	}

	lp := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
	)
	otellog.SetLoggerProvider(lp)
	p.shutdownFuncs = append(p.shutdownFuncs, lp.Shutdown)
	return nil
}

// Shutdown 按注册的逆序刷新并关闭所有 Provider，确保缓冲中的遥测数据被导出。
// 多个关闭错误会通过 errors.Join 合并返回。
func (p *Provider) Shutdown(ctx context.Context) error {
	var errs []error
	for i := len(p.shutdownFuncs) - 1; i >= 0; i-- {
		if err := p.shutdownFuncs[i](ctx); err != nil {
			errs = append(errs, err)
		}
	}
	p.shutdownFuncs = nil
	return errors.Join(errs...)
}
