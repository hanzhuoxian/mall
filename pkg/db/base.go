// Package db 封装了基于 GORM 的数据库初始化逻辑，支持 MySQL 和 PostgreSQL，
// 提供连接池配置和日志级别设置的统一入口。
package db

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// BaseOptions 是各数据库驱动共享的连接池与日志配置。
type BaseOptions struct {
	MaxIdleConnections    int            // 最大空闲连接数
	MaxOpenConnections    int            // 最大打开连接数
	MaxConnectionLifeTime time.Duration  // 连接最大存活时间，超时后自动回收
	LogLevel              int            // GORM 日志级别（1=Silent, 2=Error, 3=Warn, 4=Info）
	Logger                logger.Interface // 自定义日志实现，为 nil 时使用默认实现
}

// open 使用指定 dialector 和连接池配置打开数据库连接，是各驱动的公共创建入口。
func open(dialector gorm.Dialector, opts *BaseOptions) (*gorm.DB, error) {
	l := opts.Logger
	if l == nil {
		l = logger.Default.LogMode(logger.LogLevel(opts.LogLevel))
	}

	db, err := gorm.Open(dialector, &gorm.Config{Logger: l})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(opts.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(opts.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(opts.MaxConnectionLifeTime)

	return db, nil
}
