package cache

import "context"

type UserCache interface {
	SetUser(ctx context.Context, key string, value string) error
	GetUser(ctx context.Context, key string) (string, error)
}
