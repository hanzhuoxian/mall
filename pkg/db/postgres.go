package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresOptions struct {
	BaseOptions
	Host     string
	Port     int
	Username string
	Password string
	Database string
	SSLMode  string
	TimeZone string
}

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
