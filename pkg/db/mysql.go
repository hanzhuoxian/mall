package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// MySQLOptions 包含 MySQL 连接所需的配置，内嵌 BaseOptions 提供连接池与日志设置。
type MySQLOptions struct {
	BaseOptions
	Host     string // 地址，格式为 host:port
	Username string
	Password string
	Database string
}

// NewMySQL 使用给定配置创建并返回一个 MySQL *gorm.DB 实例。
// DSN 固定使用 utf8mb4 字符集、parseTime=true 和本地时区。
func NewMySQL(opts *MySQLOptions) (*gorm.DB, error) {
	dsn := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s`,
		opts.Username,
		opts.Password,
		opts.Host,
		opts.Database,
		true,
		"Local",
	)
	return open(mysql.Open(dsn), &opts.BaseOptions)
}
