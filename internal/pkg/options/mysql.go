package options

import (
	"time"

	"github.com/spf13/pflag"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/hanzhuoxian/mall/internal/pkg/logger"
	pkgdb "github.com/hanzhuoxian/mall/pkg/db"
)

// MySQLOptions defines options for mysql database.
type MySQLOptions struct {
	Host                  string        `json:"host,omitempty"                     mapstructure:"host"`
	Username              string        `json:"username,omitempty"                 mapstructure:"username"`
	Password              string        `json:"-"                                  mapstructure:"password"`
	Database              string        `json:"database"                           mapstructure:"database"`
	MaxIdleConnections    int           `json:"max-idle-connections,omitempty"     mapstructure:"max-idle-connections"`
	MaxOpenConnections    int           `json:"max-open-connections,omitempty"     mapstructure:"max-open-connections"`
	MaxConnectionLifeTime time.Duration `json:"max-connection-life-time,omitempty" mapstructure:"max-connection-life-time"`
	LogLevel              int           `json:"log-level"                          mapstructure:"log-level"`
}

// NewMySQLOptions 返回带有合理默认值的 MySQLOptions 实例（连接 127.0.0.1:3306）。
func NewMySQLOptions() *MySQLOptions {
	return &MySQLOptions{
		Host:                  "127.0.0.1:3306",
		Username:              "root",
		Password:              "123456",
		Database:              "",
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifeTime: time.Duration(10) * time.Second,
		LogLevel:              4,
	}
}

// Validate 校验 MySQL 选项合法性，当前无额外校验规则。
func (m *MySQLOptions) Validate() []error {
	err := []error{}
	return err
}

// AddFlags 向指定 FlagSet 注册 MySQL 连接参数的命令行 flag。
func (m *MySQLOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&m.Host, "mysql.host", m.Host, ""+
		"MySQL service host address. If left blank, the following related mysql options will be ignored.")

	fs.StringVar(&m.Username, "mysql.username", m.Username, ""+
		"Username for access to mysql service.")

	fs.StringVar(&m.Password, "mysql.password", m.Password, ""+
		"Password for access to mysql, should be used pair with password.")

	fs.StringVar(&m.Database, "mysql.database", m.Database, ""+
		"Database name for the server to use.")

	fs.IntVar(&m.MaxIdleConnections, "mysql.max-idle-connections", m.MaxIdleConnections, ""+
		"Maximum idle connections allowed to connect to mysql.")

	fs.IntVar(&m.MaxOpenConnections, "mysql.max-open-connections", m.MaxOpenConnections, ""+
		"Maximum open connections allowed to connect to mysql.")

	fs.DurationVar(&m.MaxConnectionLifeTime, "mysql.max-connection-life-time", m.MaxConnectionLifeTime, ""+
		"Maximum connection life time allowed to connect to mysql.")

	fs.IntVar(&m.LogLevel, "mysql.log-mode", m.LogLevel, ""+
		"Specify gorm log level.")
}

// NewClient 根据当前选项创建并返回一个 GORM MySQL 数据库实例。
func (o *MySQLOptions) NewClient() (*gorm.DB, error) {
	opts := &pkgdb.MySQLOptions{
		Host:     o.Host,
		Username: o.Username,
		Password: o.Password,
		Database: o.Database,
		BaseOptions: pkgdb.BaseOptions{
			MaxIdleConnections:    o.MaxIdleConnections,
			MaxOpenConnections:    o.MaxOpenConnections,
			MaxConnectionLifeTime: o.MaxConnectionLifeTime,
			LogLevel:              o.LogLevel,
			Logger:                logger.NewGormLogger(gormlogger.LogLevel(o.LogLevel)),
		},
	}

	return pkgdb.NewMySQL(opts)
}
