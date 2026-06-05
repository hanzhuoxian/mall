package log

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InfoLogger 提供 Info 级别日志能力，支持结构化字段、格式化字符串和键值对三种写法。
// Enabled 可用于在写日志前判断当前级别是否开启，避免不必要的参数构造。
type InfoLogger interface {
	Info(msg string, fields ...Field)
	Infof(fotmat string, v ...any)
	Infow(msg string, keyAndValues ...any)

	Enabled() bool
}

// Logger 是完整的日志接口，涵盖 Debug 到 Fatal 各级别，以及上下文集成、字段追加等能力。
// 所有级别均提供三种变体：强类型 Field（性能最优）、格式化字符串（fmt 风格）、键值对（sugar 风格）。
type Logger interface {
	InfoLogger
	Debug(msg string, fields ...Field)
	Debugf(format string, v ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Warn(msg string, fields ...Field)
	Warnf(format string, v ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Error(msg string, fields ...Field)
	Errorf(format string, v ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Panic(msg string, fields ...Field)
	Panicf(format string, v ...interface{})
	Panicw(msg string, keysAndValues ...interface{})
	Fatal(msg string, fields ...Field)
	Fatalf(format string, v ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})

	// V returns an InfoLogger value for a specific verbosity level.  A higher
	// verbosity level means a log message is less important.  It's illegal to
	// pass a log level less than zero.
	V(level Level) InfoLogger
	Write(p []byte) (n int, err error)

	// WithValues adds some key-value pairs of context to a logger.
	// See Info for documentation on how key/value pairs work.
	WithValues(keysAndValues ...interface{}) Logger

	// WithName adds a new element to the logger's name.
	// Successive calls with WithName continue to append
	// suffixes to the logger's name.  It's strongly recommended
	// that name segments contain only letters, digits, and hyphens
	// (see the package documentation for more information).
	WithName(name string) Logger

	// WithContext returns a copy of context in which the log value is set.
	WithContext(ctx context.Context) context.Context

	// Flush calls the underlying Core's Sync method, flushing any buffered
	// log entries. Applications should take care to call Sync before exiting.
	Flush()
}

// handleFields 将键值对参数列表（keysAndValues）转换为 zap.Field 切片。
// 若参数为奇数个或键不是 string 类型，会通过 DPanic 记录错误并提前终止转换。
func handleFields(l *zap.Logger, args []any, additional ...zap.Field) []zap.Field {
	if len(args) == 0 {
		return additional
	}

	fields := make([]zap.Field, 0, len(args)/2+len(additional))
	for i := 0; i < len(args); {
		if _, ok := args[i].(zap.Field); ok {
			l.DPanic("strongly-typed Zap Field passed to logr", zap.Any("zap field", args[i]))

			break
		}

		// make sure this isn't a mismatched key
		if i == len(args)-1 {
			l.DPanic("odd number of arguments passed as key-value pairs for logging", zap.Any("ignored key", args[i]))

			break
		}

		key, val := args[i], args[i+1]
		keyStr, isString := key.(string)
		if !isString {
			// if the key isn't a string, DPanic and stop logging
			l.DPanic(
				"non-string key argument passed to logging, ignoring all later arguments",
				zap.Any("invalid key", key),
			)

			break
		}

		fields = append(fields, zap.Any(keyStr, val))
		i += 2
	}
	return fields
}

// zapLogger 是 Logger 接口的核心实现，组合了底层 zap.Logger 和 infoLogger。
// infoLogger 负责带级别控制的 Info/Infof/Infow，zapLogger 负责其余级别及扩展方法。
type zapLogger struct {
	zapLogger *zap.Logger
	infoLogger
}

var _ Logger = &zapLogger{}

// New 根据 Options 创建一个 zapLogger 实例，同时将标准库 log 重定向至该 logger。
// opts 为 nil 时使用默认配置。
func New(opts *Options) *zapLogger {
	if opts == nil {
		opts = NewOptions()
	}

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}
	encodeLevel := zapcore.CapitalLevelEncoder
	// when output to local path, with color is forbidden
	if opts.Format == consoleFormat && opts.EnableColor {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	}

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "timestamp",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encodeLevel,
		EncodeTime:     timeEncoder,
		EncodeDuration: milliSecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	loggerConfig := &zap.Config{
		Level:             zap.NewAtomicLevelAt(zapLevel),
		Development:       opts.Development,
		DisableCaller:     opts.DisableCaller,
		DisableStacktrace: opts.DisableStacktrace,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         opts.Format,
		EncoderConfig:    encoderConfig,
		OutputPaths:      opts.OutputPaths,
		ErrorOutputPaths: opts.ErrorOutputPaths,
	}

	var err error
	l, err := loggerConfig.Build(zap.AddStacktrace(zapcore.PanicLevel), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	logger := &zapLogger{
		zapLogger: l.Named(opts.Name),
		infoLogger: infoLogger{
			log:   l,
			level: zap.InfoLevel,
		},
	}
	zap.RedirectStdLog(l)

	return logger
}

func (l *zapLogger) V(level Level) InfoLogger {
	if l.zapLogger.Core().Enabled(level) {
		return &infoLogger{
			level: level,
			log:   l.zapLogger,
		}
	}

	return disabledInfoLogger
}

// Write 实现 io.Writer 接口，将写入内容以 Info 级别记录，便于与标准库日志集成。
func (l *zapLogger) Write(p []byte) (n int, err error) {
	l.zapLogger.Info(string(p))

	return len(p), nil
}

func (l *zapLogger) WithValues(keysAndValues ...interface{}) Logger {
	newLogger := l.zapLogger.With(handleFields(l.zapLogger, keysAndValues)...)

	return NewLogger(newLogger)
}

func WithName(s string) Logger {
	return std.WithName(s)
}

func (l *zapLogger) WithName(name string) Logger {
	return NewLogger(l.zapLogger.Named(name))
}

func (l *zapLogger) Flush() {
	l.zapLogger.Sync()
}

// NewLogger creates a new logr.Logger using the given Zap Logger to log.
func NewLogger(l *zap.Logger) Logger {
	return &zapLogger{
		zapLogger: l,
		infoLogger: infoLogger{
			log:   l,
			level: zap.InfoLevel,
		},
	}
}

// CheckIntLevel used for other log wrapper such as klog which return if logging a
// message at the specified level is enabled.
func CheckIntLevel(level int32) bool {
	var lvl zapcore.Level
	if level < 5 {
		lvl = zapcore.InfoLevel
	} else {
		lvl = zapcore.DebugLevel
	}
	checkEntry := std.zapLogger.Check(lvl, "")

	return checkEntry != nil
}

// Debug method output debug level log.
func Debug(msg string, fields ...Field) {
	std.zapLogger.Debug(msg, fields...)
}

func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.zapLogger.Debug(msg, fields...)
}

// Debugf method output debug level log.
func Debugf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Debugf(format, v...)
}

func (l *zapLogger) Debugf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Debugf(format, v...)
}

// Debugw method output debug level log.
func Debugw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Debugw(msg, keysAndValues...)
}

func (l *zapLogger) Debugw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Debugw(msg, keysAndValues...)
}

// Info method output info level log.
func Info(msg string, fields ...Field) {
	std.zapLogger.Info(msg, fields...)
}

func (l *zapLogger) Info(msg string, fields ...Field) {
	l.zapLogger.Info(msg, fields...)
}

// Infof method output info level log.
func Infof(format string, v ...interface{}) {
	std.zapLogger.Sugar().Infof(format, v...)
}

func (l *zapLogger) Infof(format string, v ...interface{}) {
	l.zapLogger.Sugar().Infof(format, v...)
}

// Infow method output info level log.
func Infow(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Infow(msg, keysAndValues...)
}

func (l *zapLogger) Infow(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Infow(msg, keysAndValues...)
}

// Warn method output warning level log.
func Warn(msg string, fields ...Field) {
	std.zapLogger.Warn(msg, fields...)
}

func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.zapLogger.Warn(msg, fields...)
}

// Warnf method output warning level log.
func Warnf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Warnf(format, v...)
}

func (l *zapLogger) Warnf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Warnf(format, v...)
}

// Warnw method output warning level log.
func Warnw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Warnw(msg, keysAndValues...)
}

func (l *zapLogger) Warnw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Warnw(msg, keysAndValues...)
}

// Error method output error level log.
func Error(msg string, fields ...Field) {
	std.zapLogger.Error(msg, fields...)
}

func (l *zapLogger) Error(msg string, fields ...Field) {
	l.zapLogger.Error(msg, fields...)
}

// Errorf method output error level log.
func Errorf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Errorf(format, v...)
}

func (l *zapLogger) Errorf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Errorf(format, v...)
}

// Errorw method output error level log.
func Errorw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Errorw(msg, keysAndValues...)
}

func (l *zapLogger) Errorw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Errorw(msg, keysAndValues...)
}

// Panic method output panic level log and shutdown application.
func Panic(msg string, fields ...Field) {
	std.zapLogger.Panic(msg, fields...)
}

func (l *zapLogger) Panic(msg string, fields ...Field) {
	l.zapLogger.Panic(msg, fields...)
}

// Panicf method output panic level log and shutdown application.
func Panicf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Panicf(format, v...)
}

func (l *zapLogger) Panicf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Panicf(format, v...)
}

// Panicw method output panic level log.
func Panicw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Panicw(msg, keysAndValues...)
}

func (l *zapLogger) Panicw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Panicw(msg, keysAndValues...)
}

// Fatal method output fatal level log.
func Fatal(msg string, fields ...Field) {
	std.zapLogger.Fatal(msg, fields...)
}

func (l *zapLogger) Fatal(msg string, fields ...Field) {
	l.zapLogger.Fatal(msg, fields...)
}

// Fatalf method output fatal level log.
func Fatalf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Fatalf(format, v...)
}

func (l *zapLogger) Fatalf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Fatalf(format, v...)
}

// Fatalw method output Fatalw level log.
func Fatalw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Fatalw(msg, keysAndValues...)
}

func (l *zapLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Fatalw(msg, keysAndValues...)
}

// L 从 context 中提取预设字段（requestID、username、watcherName）并附加到 logger，
// 便于在链路追踪场景下自动携带请求元信息。
func L(ctx context.Context) *zapLogger {
	return std.L(ctx)
}

func (l *zapLogger) L(ctx context.Context) *zapLogger {
	lg := l.clone()
	fields := []Field{}
	if requestID := ctx.Value(KeyRequestID); requestID != nil {
		fields = append(fields, zap.Any(KeyRequestID, requestID))
	}
	if username := ctx.Value(KeyUsername); username != nil {
		fields = append(fields, zap.Any(KeyUsername, username))
	}
	if watcherName := ctx.Value(KeyWatcherName); watcherName != nil {
		fields = append(fields, zap.Any(KeyWatcherName, watcherName))
	}

	lg.zapLogger = lg.zapLogger.With(fields...)

	return lg
}

//nolint:predeclared
// clone 浅拷贝当前 logger，用于派生带额外字段的子 logger，避免修改原实例。
func (l *zapLogger) clone() *zapLogger {
	copy := *l

	return &copy
}
