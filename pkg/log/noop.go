package log

// noopInfoLogger 是空操作的 InfoLogger 实现，所有方法均为无操作。
// 当请求的日志级别未开启时，V() 返回此实例以避免无效写入。
type noopInfoLogger struct{}

// Enabled 始终返回 false，告知调用方该级别已被禁用。
func (*noopInfoLogger) Enabled() bool {
	return false
}

func (l *noopInfoLogger) Info(_ string, _ ...Field) {}
func (l *noopInfoLogger) Infof(_ string, _ ...any)  {}
func (l *noopInfoLogger) Infow(_ string, _ ...any)  {}

// disabledInfoLogger 是全局共享的空 InfoLogger 单例，供 V() 在级别未开启时返回。
var disabledInfoLogger = &noopInfoLogger{}
