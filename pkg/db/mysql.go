package db

import (
	"fmt"
	"net/url"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// MySQLOptions 包含 MySQL 连接所需的配置，内嵌 BaseOptions 提供连接池与日志设置。
type MySQLOptions struct {
	Host     string   // 主库地址，格式为 host:port
	Replicas []string // 从库地址列表，为空时不启用读写分离
	Username string
	Password string
	Database string
	Charset  string // 字符集，默认 utf8mb4
	Loc      string // 时区，默认 Asia/Shanghai
	BaseOptions
}

// NewMySQL 使用给定配置创建并返回一个 MySQL *gorm.DB 实例。
// 若 Replicas 非空，自动注册 dbresolver 插件实现读写分离：SELECT 路由从库，写操作路由主库。
func NewMySQL(opts *MySQLOptions) (*gorm.DB, error) {
	dsn := mysqlDSN(opts.Host, opts)
	db, err := open(mysql.Open(dsn), &opts.BaseOptions)
	if err != nil {
		return nil, err
	}
	if len(opts.Replicas) == 0 {
		return db, nil
	}

	replicas := make([]gorm.Dialector, 0, len(opts.Replicas))
	for _, replica := range opts.Replicas {
		replicas = append(replicas, mysql.Open(mysqlDSN(replica, opts)))
	}

	if err := db.Use(dbresolver.Register(dbresolver.Config{
		Replicas: replicas,
		Policy:   dbresolver.RandomPolicy{},
	})); err != nil {
		return nil, fmt.Errorf("register dbresolver: %w", err)
	}
	return db, nil
}

func mysqlDSN(host string, opts *MySQLOptions) string {
	charset := opts.Charset
	if charset == "" {
		charset = "utf8mb4"
	}
	loc := opts.Loc
	if loc == "" {
		loc = "Local"
	}
	return fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=%s`,
		opts.Username,
		opts.Password,
		host,
		opts.Database,
		charset,
		true,
		url.QueryEscape(loc),
	)
}
