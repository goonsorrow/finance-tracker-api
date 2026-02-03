package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type AuthRedis struct {
	rdb *redis.Client
}

func NewAuthRedis(rdb *redis.Client) *AuthRedis {
	return &AuthRedis{
		rdb: rdb,
	}
}

func (c AuthRedis) CacheRefreshToken(ctx context.Context, key string, refreshTTL time.Duration) error {
	return c.rdb.Set(ctx, key, "valid", refreshTTL).Err()
}

func (c AuthRedis) CheckRefreshToken(ctx context.Context, key string) (int, error) {
	exists, err := c.rdb.Exists(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return int(exists), nil
}

func (c AuthRedis) DeleteRefreshToken(ctx context.Context, key string) error {
	return c.rdb.Del(ctx, key).Err()
}
