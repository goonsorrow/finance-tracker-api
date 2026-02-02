package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Authorization interface {
	CacheRefreshSession(ctx context.Context, key string, refreshTTL time.Duration) error
	CheckRefreshToken(ctx context.Context, key string) (int, error)
	//DeleteRefreshSession
}

// type Category interface{

// }

type Cache struct {
	Authorization
	// Category
}

func NewCache(rdb *redis.Client) *Cache {
	return &Cache{
		Authorization: NewAuthRedis(rdb),
	}
}
