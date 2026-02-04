package cache

import (
	"context"
	"fmt"
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

func (c AuthRedis) GetUserRefreshSessions(ctx context.Context, userId int) ([]string, error) {
	match := fmt.Sprintf("refresh:userId:%d:*", userId)

	var keys []string
	i := c.rdb.Scan(ctx, 0, match, 0).Iterator()
	for i.Next(ctx) {
		keys = append(keys, i.Val())
	}
	if err := i.Err(); err != nil {
		return nil, err
	}
	return keys, nil
}

func (c AuthRedis) DeleteRefreshToken(ctx context.Context, key string) error {
	return c.rdb.Del(ctx, key).Err()
}

func (c AuthRedis) DeleteAllRefreshTokens(ctx context.Context, keys []string) error {
	return c.rdb.Del(ctx, keys...).Err()
}
