package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// PostgresOptions 包含 PostgreSQL 连接所需的配置，内嵌 BaseOptions 提供连接池与日志设置。
type PostgresOptions struct {
	Host     string   // 主库地址
	Replicas []string // 从库地址列表，为空时不启用读写分离
	Username string
	Password string
	Database string
	SSLMode  string // SSL 模式，默认 "disable"
	TimeZone string // 时区，默认 "Asia/Shanghai"
	BaseOptions
	Port int
}

func postgresDSN(host string, opts *PostgresOptions) string {
	sslMode := opts.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}
	timeZone := opts.TimeZone
	if timeZone == "" {
		timeZone = "Asia/Shanghai"
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		host, opts.Port, opts.Username, opts.Password, opts.Database, sslMode, timeZone,
	)
}

// NewPostgres 使用给定配置创建并返回一个 PostgreSQL *gorm.DB 实例。
// SSLMode 和 TimeZone 未设置时使用安全的默认值。
// 若 Replicas 非空，自动注册 dbresolver 插件实现读写分离：SELECT 路由从库，写操作路由主库。
func NewPostgres(opts *PostgresOptions) (*gorm.DB, error) {
	db, err := open(postgres.Open(postgresDSN(opts.Host, opts)), &opts.BaseOptions)
	if err != nil {
		return nil, err
	}
	if len(opts.Replicas) == 0 {
		return db, nil
	}

	replicas := make([]gorm.Dialector, 0, len(opts.Replicas))
	for _, host := range opts.Replicas {
		replicas = append(replicas, postgres.Open(postgresDSN(host, opts)))
	}
	if err := db.Use(dbresolver.Register(dbresolver.Config{
		Replicas: replicas,
		Policy:   dbresolver.RandomPolicy{},
	})); err != nil {
		return nil, fmt.Errorf("register dbresolver: %w", err)
	}
	return db, nil
}
