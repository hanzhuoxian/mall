package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLOptions struct {
	BaseOptions
	Host     string
	Username string
	Password string
	Database string
}

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
