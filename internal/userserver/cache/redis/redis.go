package redis

import (
	"fmt"

	"github.com/google/wire"
	"github.com/hanzhuoxian/mall/internal/pkg/options"
	"github.com/hanzhuoxian/mall/internal/userserver/cache"
	goredis "github.com/redis/go-redis/v9"
)

// ProviderSet is used by Wire.
var ProviderSet = wire.NewSet(NewCachestore)

type cachestore struct {
	client goredis.UniversalClient
}

func (c *cachestore) User() cache.UserCache {
	return &userCache{client: c.client}
}

func (c *cachestore) Close() error {
	return c.client.Close()
}

// NewCachestore creates a Redis-backed cache.Factory.
func NewCachestore(opts *options.RedisOptions) (cache.Factory, error) {
	r, err := opts.NewClient()
	if err != nil {
		return nil, fmt.Errorf("connect redis: %w", err)
	}
	return &cachestore{client: r}, nil
}
