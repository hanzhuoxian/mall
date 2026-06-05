// Package storage 提供 Redis 客户端的初始化封装，基于 go-redis/v9 的 UniversalClient，
// 支持单节点、Sentinel 和 Cluster 三种模式（通过 MasterName 和 Addrs 数量自动选择）。
package storage

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisOptions 包含 Redis 连接所需的配置项。
// MasterName 非空时使用 Sentinel 模式；Addrs 多于一个时使用 Cluster 模式；否则为单节点模式。
type RedisOptions struct {
	Addrs        []string      // 节点地址列表，格式为 host:port
	Username     string
	Password     string
	DB           int           // 数据库编号，Cluster 模式下无效
	DialTimeout  time.Duration // 建立连接的超时时间
	ReadTimeout  time.Duration // 读操作超时时间
	WriteTimeout time.Duration // 写操作超时时间
	PoolSize     int           // 每个节点的最大连接数
	MinIdleConns int           // 每个节点维持的最小空闲连接数
	MasterName   string        // Sentinel 主节点名称，非空时启用 Sentinel 模式
}

// NewRedis 创建并返回一个 Redis 通用客户端，初始化后会发送 PING 验证连通性。
// 5 秒内无法连通则返回错误。
func NewRedis(opts *RedisOptions) (redis.UniversalClient, error) {
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:        opts.Addrs,
		Username:     opts.Username,
		Password:     opts.Password,
		DB:           opts.DB,
		DialTimeout:  opts.DialTimeout,
		ReadTimeout:  opts.ReadTimeout,
		WriteTimeout: opts.WriteTimeout,
		PoolSize:     opts.PoolSize,
		MinIdleConns: opts.MinIdleConns,
		MasterName:   opts.MasterName,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
