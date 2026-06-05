package log

import (
	"time"

	"go.uber.org/zap/zapcore"
)

// timeEncoder 将时间格式化为 "2006-01-02 15:04:05.000"，精度到毫秒，便于人工阅读。
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// milliSecondsDurationEncoder 将 Duration 编码为毫秒浮点数，避免纳秒数值过大难以阅读。
func milliSecondsDurationEncoder(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendFloat64(float64(d) / float64(time.Millisecond))
}
