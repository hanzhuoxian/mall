package cache

import "context"

// UserCache 定义了用户缓存的读写操作接口。
type UserCache interface {
	// SetUser 将用户信息以 key-value 形式写入缓存。
	SetUser(ctx context.Context, key string, value string) error
	// GetUser 从缓存中读取指定 key 的用户信息。
	GetUser(ctx context.Context, key string) (string, error)
}
