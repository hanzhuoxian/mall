package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgresOptions 包含 PostgreSQL 连接所需的配置，内嵌 BaseOptions 提供连接池与日志设置。
type PostgresOptions struct {
	Host     string
	Username string
	Password string
	Database string
	SSLMode  string // SSL 模式，默认 "disable"
	TimeZone string // 时区，默认 "Asia/Shanghai"
	BaseOptions
	Port int
}

// NewPostgres 使用给定配置创建并返回一个 PostgreSQL *gorm.DB 实例。
// SSLMode 和 TimeZone 未设置时使用安全的默认值。
func NewPostgres(opts *PostgresOptions) (*gorm.DB, error) {
	sslMode := opts.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}
	timeZone := opts.TimeZone
	if timeZone == "" {
		timeZone = "Asia/Shanghai"
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		opts.Host,
		opts.Port,
		opts.Username,
		opts.Password,
		opts.Database,
		sslMode,
		timeZone,
	)
	return open(postgres.Open(dsn), &opts.BaseOptions)
}
