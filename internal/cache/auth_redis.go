package cache

import (
	"context"
	"errors"
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

func (a AuthRedis) CacheRefreshSession(ctx context.Context, key string, refreshTTL time.Duration) error {
	return a.rdb.Set(ctx, key, "valid", refreshTTL).Err()
}

func (a AuthRedis) CheckRefreshToken(ctx context.Context, key string) (int, error) {
	exists, _ := a.rdb.Exists(ctx, key).Result()

	if exists == 0 {
		return int(exists), errors.New("token expired/revoked")
	} else {
		return int(exists), errors.New("token still valid")
	}
}

// func deleteRefreshSession
