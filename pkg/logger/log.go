package logger

import (
	"errors"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelzap"
	otellog "go.opentelemetry.io/otel/log/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	samplingInterval   = time.Second
	samplingInitial    = 100
	samplingThereafter = 100
)

// initMu serializes installation of the zap and standard-library globals.
// Log calls read globals through zap's concurrency-safe accessors.
var initMu sync.Mutex

// New creates a native zap logger without changing global loggers.
// Both local output and OTel use the configured level and sampling policy.
func New(opts *Options) (*zap.Logger, error) {
	if opts == nil {
		opts = NewOptions()
	}
	if err := errors.Join(opts.Validate()...); err != nil {
		return nil, err
	}
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(opts.Level)); err != nil {
		return nil, err
	}
	format := strings.ToLower(opts.Format)
	encodeLevel := zapcore.CapitalLevelEncoder
	if format == consoleFormat && opts.EnableColor {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	}
	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       opts.Development,
		DisableCaller:     opts.DisableCaller,
		DisableStacktrace: opts.DisableStacktrace,
		Encoding:          format,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message", LevelKey: "level", TimeKey: "timestamp",
			NameKey: "logger", CallerKey: "caller", StacktraceKey: "stacktrace",
			LineEnding:  zapcore.DefaultLineEnding,
			EncodeLevel: encodeLevel, EncodeTime: timeEncoder,
			EncodeDuration: milliSecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
		OutputPaths:      opts.OutputPaths,
		ErrorOutputPaths: opts.ErrorOutputPaths,
	}
	var options []zap.Option
	if !opts.DisableStacktrace {
		options = append(options, zap.AddStacktrace(zapcore.PanicLevel))
	}
	l, err := cfg.Build(options...)
	if err != nil {
		return nil, err
	}
	otelCore := otelzap.NewCore(opts.Name, otelzap.WithLoggerProvider(otellog.GetLoggerProvider()))
	// The local core supports cfg.Level, even if the OTel provider is disabled.
	core, err := zapcore.NewIncreaseLevelCore(zapcore.NewTee(l.Core(), otelCore), cfg.Level)
	if err != nil {
		return nil, err
	}
	core = zapcore.NewSamplerWithOptions(core, samplingInterval, samplingInitial, samplingThereafter)
	return l.WithOptions(zap.WrapCore(func(zapcore.Core) zapcore.Core { return core })).Named(opts.Name), nil
}

// Init installs the logger for zap.L(), zap.S() and the standard-library log.
// Initialize once at startup, after telemetry, and Sync the installed logger on exit.
// Invalid configuration leaves the existing global logger unchanged.
func Init(opts *Options) error {
	initMu.Lock()
	defer initMu.Unlock()
	l, err := New(opts)
	if err != nil {
		return err
	}
	zap.ReplaceGlobals(l)
	zap.RedirectStdLog(l)
	return nil
}
