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

func GetCacheFactoryOr(opts *options.RedisOptions) (f cache.Factory, err error) {
	if opts == nil && cacheFactory == nil {
		return nil, errors.New("get redis factory failed")
	}
	once.Do(func() {
		var r redis.UniversalClient
		r, err = opts.NewClient()
		cacheFactory = &cachestore{Client: r}
	})

	if cacheFactory == nil || err != nil {
		return nil, errors.New("create redis factory failed")
	}

	return cacheFactory, nil
}
