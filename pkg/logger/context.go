package logger

import "context"

// key 是 context 存储 logger 的私有键类型，防止与其他包的 key 冲突。
type key int

const (
	logContextKey key = iota // logContextKey 是存储 Logger 实例的 context key。
)

// WithContext 将当前 logger 存入 context 并返回新 context，便于跨层传递。
func (l *zapLogger) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, logContextKey, l)
}

// WithContext 将全局 logger 存入 context，是包级别的快捷方式。
func WithContext(ctx context.Context) context.Context {
	return std.WithContext(ctx)
}

// FromContext 从 context 中取出 Logger；若未设置，返回名为 "Unknown-Context" 的 logger
// 以保证调用方始终拿到有效的 Logger，不会触发 nil panic。
func FromContext(ctx context.Context) Logger {
	if ctx != nil {
		logger := ctx.Value(logContextKey)
		if logger != nil {
			return logger.(Logger)
		}
	}

	return WithName("Unknown-Context")
}
