// Package captcha provides a Redis-backed store for base64Captcha.
package captcha

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	keyPrefix     = "captcha:"
	defaultExpiry = 5 * time.Minute
)

// RedisStore implements base64Captcha.Store backed by Redis.
type RedisStore struct {
	rdb    redis.UniversalClient
	expiry time.Duration
}

// NewRedisStore creates a RedisStore with the given client and a 5-minute expiry.
func NewRedisStore(rdb redis.UniversalClient) *RedisStore {
	return &RedisStore{rdb: rdb, expiry: defaultExpiry}
}

func (s *RedisStore) Set(id string, value string) error {
	return s.rdb.Set(context.Background(), keyPrefix+id, value, s.expiry).Err()
}

func (s *RedisStore) Get(id string, clear bool) string {
	key := keyPrefix + id
	val, err := s.rdb.Get(context.Background(), key).Result()
	if err != nil {
		return ""
	}
	if clear {
		s.rdb.Del(context.Background(), key)
	}
	return val
}

func (s *RedisStore) Verify(id, answer string, clear bool) bool {
	return s.Get(id, clear) == answer
}
