package logger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"

	"github.com/hanzhuoxian/mall/pkg/log"
)

type GormLogger struct {
	logger                    *zap.Logger
	logLevel                  gormlogger.LogLevel
	slowThreshold             time.Duration
	ignoreRecordNotFoundError bool
}

var _ gormlogger.Interface = (*GormLogger)(nil)

func NewGormLogger(level gormlogger.LogLevel) *GormLogger {
	return &GormLogger{
		logger:                    log.ZapLogger(),
		logLevel:                  gormlogger.LogLevel(level),
		slowThreshold:             200 * time.Millisecond,
		ignoreRecordNotFoundError: true,
	}
}

func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	cp := *l
	cp.logLevel = level
	return &cp
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...any) {
	if l.logLevel >= gormlogger.Info {
		l.logger.Sugar().Infof("[gorm] "+msg, data...)
	}
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...any) {
	if l.logLevel >= gormlogger.Warn {
		l.logger.Sugar().Warnf("[gorm] "+msg, data...)
	}
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...any) {
	if l.logLevel >= gormlogger.Error {
		l.logger.Sugar().Errorf("[gorm] "+msg, data...)
	}
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	caller := utils.FileWithLineNum()

	switch {
	case err != nil && l.logLevel >= gormlogger.Error &&
		(!errors.Is(err, gormlogger.ErrRecordNotFound) || !l.ignoreRecordNotFoundError):
		sql, rows := fc()
		l.logger.Error("[gorm] trace",
			zap.String("caller", caller),
			zap.Error(err),
			zap.String("elapsed", fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	case elapsed > l.slowThreshold && l.slowThreshold > 0 && l.logLevel >= gormlogger.Warn:
		sql, rows := fc()
		l.logger.Warn("[gorm] slow sql",
			zap.String("caller", caller),
			zap.String("threshold", l.slowThreshold.String()),
			zap.String("elapsed", fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	case l.logLevel >= gormlogger.Info:
		sql, rows := fc()
		l.logger.Info("[gorm] trace",
			zap.String("caller", caller),
			zap.String("elapsed", fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	}
}
