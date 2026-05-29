package options

import (
	"time"

	"github.com/hanzhuoxian/mall/pkg/storage"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/pflag"
)

type RedisOptions struct {
	Addrs        []string      `json:"addrs,omitempty"          mapstructure:"addrs"`
	Username     string        `json:"username,omitempty"       mapstructure:"username"`
	Password     string        `json:"-"                        mapstructure:"password"`
	DB           int           `json:"db,omitempty"             mapstructure:"db"`
	DialTimeout  time.Duration `json:"dial-timeout,omitempty"   mapstructure:"dial-timeout"`
	ReadTimeout  time.Duration `json:"read-timeout,omitempty"   mapstructure:"read-timeout"`
	WriteTimeout time.Duration `json:"write-timeout,omitempty"  mapstructure:"write-timeout"`
	PoolSize     int           `json:"pool-size,omitempty"      mapstructure:"pool-size"`
	MinIdleConns int           `json:"min-idle-conns,omitempty" mapstructure:"min-idle-conns"`
	MasterName   string        `json:"master-name,omitempty"    mapstructure:"master-name"`
}

func NewRedisOptions() *RedisOptions {
	return &RedisOptions{
		Addrs:        []string{"127.0.0.1:6379"},
		DB:           0,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 2,
	}
}

func (o *RedisOptions) Validate() []error {
	return []error{}
}

func (o *RedisOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringSliceVar(&o.Addrs, "redis.addrs", o.Addrs,
		"Comma-separated list of Redis server addresses (host:port). Multiple addresses enable cluster mode.")

	fs.StringVar(&o.Username, "redis.username", o.Username,
		"Username for Redis authentication (Redis 6+ ACL).")

	fs.StringVar(&o.Password, "redis.password", o.Password,
		"Password for Redis authentication.")

	fs.IntVar(&o.DB, "redis.db", o.DB,
		"Database index to select after connecting (only for standalone mode).")

	fs.DurationVar(&o.DialTimeout, "redis.dial-timeout", o.DialTimeout,
		"Timeout for establishing new connections to Redis.")

	fs.DurationVar(&o.ReadTimeout, "redis.read-timeout", o.ReadTimeout,
		"Timeout for socket reads.")

	fs.DurationVar(&o.WriteTimeout, "redis.write-timeout", o.WriteTimeout,
		"Timeout for socket writes.")

	fs.IntVar(&o.PoolSize, "redis.pool-size", o.PoolSize,
		"Maximum number of socket connections in the pool.")

	fs.IntVar(&o.MinIdleConns, "redis.min-idle-conns", o.MinIdleConns,
		"Minimum number of idle connections in the pool.")

	fs.StringVar(&o.MasterName, "redis.master-name", o.MasterName,
		"Sentinel master name. When set, enables Sentinel mode.")
}

func (o *RedisOptions) NewClient() (redis.UniversalClient, error) {
	return storage.NewRedis(&storage.RedisOptions{
		Addrs:        o.Addrs,
		Username:     o.Username,
		Password:     o.Password,
		DB:           o.DB,
		DialTimeout:  o.DialTimeout,
		ReadTimeout:  o.ReadTimeout,
		WriteTimeout: o.WriteTimeout,
		PoolSize:     o.PoolSize,
		MinIdleConns: o.MinIdleConns,
		MasterName:   o.MasterName,
	})
}
