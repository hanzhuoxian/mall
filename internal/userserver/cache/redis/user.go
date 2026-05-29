package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type userCache struct {
	client redis.UniversalClient
}

func (u *userCache) SetUser(ctx context.Context, key string, value string) error {
	return u.client.Set(ctx, key, value, 0).Err()
}

func (u *userCache) GetUser(ctx context.Context, key string) (string, error) {
	return u.client.Get(ctx, key).Result()
}
