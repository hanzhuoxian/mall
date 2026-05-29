package redis

import (
	"errors"
	"sync"

	"github.com/hanzhuoxian/mall/internal/pkg/options"
	"github.com/hanzhuoxian/mall/internal/userserver/cache"
	"github.com/redis/go-redis/v9"
)

type cachestore struct {
	Client redis.UniversalClient
}

var (
	cacheFactory cache.Factory
	once         sync.Once
)

func (c *cachestore) User() cache.UserCache {
	return &userCache{client: c.Client}
}

func (c *cachestore) Close() error {
	return c.Client.Close()
}

func GetCacheFactoryOr(opts *options.RedisOptions) (f cache.Factory, err error) {
	if opts == nil && cacheFactory == nil {
		return nil, errors.New("get redis factory failed")
	}
	once.Do(func() {
		var r redis.UniversalClient
		r, err = opts.NewClient()
		if err != nil {
			return
		}
		cacheFactory = &cachestore{Client: r}
	})

	if err != nil {
		return nil, err
	}
	if cacheFactory == nil {
		return nil, errors.New("create redis factory failed")
	}

	return cacheFactory, nil
}
