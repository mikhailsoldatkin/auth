package cache

import (
	"context"
	"time"
)

type RedisClient interface {
	HashSet(ctx context.Context, key string, values interface{}) error
	Set(ctx context.Context, key string, value interface{}) error
	HGetAll(ctx context.Context, key string) ([]interface{}, error)
	Get(ctx context.Context, key string) (interface{}, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Ping(ctx context.Context) error
	Keys(ctx context.Context, s string) ([]string, error)
}
