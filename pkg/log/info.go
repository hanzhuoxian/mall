package log

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// infoLogger 是带可配置级别的 InfoLogger 实现，被 zapLogger 内嵌用于支持 V() 按级别过滤。
// 通过 zap.Logger.Check 预检查级别，在日志未开启时跳过格式化，减少无效开销。
type infoLogger struct {
	log   *zap.Logger
	level zapcore.Level
}

// Enabled 始终返回 true，表示该 infoLogger 对应的级别已开启（禁用时使用 noopInfoLogger）。
func (l *infoLogger) Enabled() bool {
	return true
}

// Info 以强类型 Field 写入日志，级别未开启时无任何开销。
func (l *infoLogger) Info(msg string, fields ...Field) {
	if checkedEntry := l.log.Check(l.level, msg); checkedEntry != nil {
		checkedEntry.Write(fields...)
	}
}

// Infof 以 fmt 格式化字符串写入日志。
func (l *infoLogger) Infof(format string, args ...any) {
	if checkedEntry := l.log.Check(l.level, fmt.Sprintf(format, args...)); checkedEntry != nil {
		checkedEntry.Write()
	}
}

// Infow 以键值对形式写入日志，键必须为 string 类型。
func (l *infoLogger) Infow(msg string, keysAndValues ...any) {
	if checkedEntry := l.log.Check(l.level, msg); checkedEntry != nil {
		checkedEntry.Write(handleFields(l.log, keysAndValues)...)
	}
}
