package log

import (
	"log"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	std = New(NewOptions()) // std 是包级别的全局默认 logger。
	mu  sync.Mutex          // mu 保护 std 的并发替换。
)

// Init 用新配置替换全局 logger，线程安全。通常在应用启动时调用一次。
func Init(opts *Options) {
	mu.Lock()
	defer mu.Unlock()
	std = New(opts)
}

// SugaredLogger returns global sugared logger.
func SugaredLogger() *zap.SugaredLogger {
	return std.zapLogger.Sugar()
}

// StdErrLogger returns logger of standard library which writes to supplied zap
// logger at error level.
func StdErrLogger() *log.Logger {
	if std == nil {
		return nil
	}
	if l, err := zap.NewStdLogAt(std.zapLogger, zapcore.ErrorLevel); err == nil {
		return l
	}

	return nil
}

// StdInfoLogger returns logger of standard library which writes to supplied zap
// logger at info level.
func StdInfoLogger() *log.Logger {
	if std == nil {
		return nil
	}
	if l, err := zap.NewStdLogAt(std.zapLogger, zapcore.InfoLevel); err == nil {
		return l
	}

	return nil
}

// Flush 刷新全局 logger 的缓冲区，应用退出前应调用以防日志丢失。
func Flush() {
	std.Flush()
}

// V 返回指定级别的全局 InfoLogger；级别未开启时返回空操作 logger。
func V(level Level) InfoLogger {
	return std.V(level)
}

// WithValues creates a child logger and adds Zap fields to it.
func WithValues(keysAndValues ...any) Logger { return std.WithValues(keysAndValues...) }

// ZapLogger 返回全局 logger 底层的 *zap.Logger，供需要直接使用 zap API 的场景使用。
func ZapLogger() *zap.Logger {
	return std.zapLogger
}
