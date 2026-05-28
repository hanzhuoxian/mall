package db

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// BaseOptions holds connection pool and logging settings shared by all drivers.
type BaseOptions struct {
	MaxIdleConnections    int
	MaxOpenConnections    int
	MaxConnectionLifeTime time.Duration
	LogLevel              int
	Logger                logger.Interface
}

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
