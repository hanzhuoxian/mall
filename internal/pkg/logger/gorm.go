// Package logger 提供适配 GORM 日志接口（gormlogger.Interface）的 zap 实现，
// 支持慢查询告警和 SQL 执行链路追踪。
package logger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"

	"github.com/hanzhuoxian/mall/pkg/logger"
)

// GormLogger 是基于 zap 的 GORM 日志实现，支持日志级别控制、慢查询阈值告警
// 以及可选的 RecordNotFound 错误忽略。
type GormLogger struct {
	logger                    *zap.Logger
	logLevel                  gormlogger.LogLevel
	slowThreshold             time.Duration // 慢查询阈值，超过此时间的 SQL 会以 Warn 级别输出
	ignoreRecordNotFoundError bool          // 为 true 时不将 RecordNotFound 记录为错误日志
}

var _ gormlogger.Interface = (*GormLogger)(nil)

// NewGormLogger 使用注入的 zap logger，慢查询阈值默认为 200ms，忽略 RecordNotFound 错误。
func NewGormLogger(l *zap.Logger, level gormlogger.LogLevel) *GormLogger {
	return &GormLogger{
		logger:                    l,
		logLevel:                  level,
		slowThreshold:             200 * time.Millisecond,
		ignoreRecordNotFoundError: true,
	}
}

// LogMode 返回一个调整了日志级别的新 GormLogger 副本，实现 gormlogger.Interface。
func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	cp := *l
	cp.logLevel = level
	return &cp
}

// Info 在日志级别 >= Info 时输出 info 级别的 GORM 日志，实现 gormlogger.Interface。
func (l *GormLogger) Info(ctx context.Context, msg string, data ...any) {
	if l.logLevel >= gormlogger.Info {
		logger.WithContext(ctx, l.logger).Sugar().Infof("[gorm] "+msg, data...)
	}
}

// Warn 在日志级别 >= Warn 时输出 warn 级别的 GORM 日志，实现 gormlogger.Interface。
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...any) {
	if l.logLevel >= gormlogger.Warn {
		logger.WithContext(ctx, l.logger).Sugar().Warnf("[gorm] "+msg, data...)
	}
}

// Error 在日志级别 >= Error 时输出 error 级别的 GORM 日志，实现 gormlogger.Interface。
func (l *GormLogger) Error(ctx context.Context, msg string, data ...any) {
	if l.logLevel >= gormlogger.Error {
		logger.WithContext(ctx, l.logger).Sugar().Errorf("[gorm] "+msg, data...)
	}
}

// Trace 记录每条 SQL 的执行详情：出错时记录 error，超过慢查询阈值时记录 warn，其余记录 info。
// 实现 gormlogger.Interface。
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel <= gormlogger.Silent {
		return
	}

	log := logger.WithContext(ctx, l.logger)
	elapsed := time.Since(begin)
	// CallerFrame skips this adapter before filtering GORM's internal frames.
	// FileWithLineNum adds another frame and would point back into this file.
	frame := utils.CallerFrame()
	caller := ""
	if frame.PC != 0 {
		caller = fmt.Sprintf("%s:%d", frame.File, frame.Line)
	}

	switch {
	case err != nil && l.logLevel >= gormlogger.Error &&
		(!errors.Is(err, gormlogger.ErrRecordNotFound) || !l.ignoreRecordNotFoundError):
		sql, rows := fc()
		log.Error("[gorm] trace",
			zap.String("sql_caller", caller),
			zap.Error(err),
			zap.String("elapsed", fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	case elapsed > l.slowThreshold && l.slowThreshold > 0 && l.logLevel >= gormlogger.Warn:
		sql, rows := fc()
		log.Warn("[gorm] slow sql",
			zap.String("sql_caller", caller),
			zap.String("threshold", l.slowThreshold.String()),
			zap.String("elapsed", fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	case l.logLevel >= gormlogger.Info:
		sql, rows := fc()
		log.Info("[gorm] trace",
			zap.String("sql_caller", caller),
			zap.String("elapsed", fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	}
}
