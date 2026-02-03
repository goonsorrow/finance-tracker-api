package cache

import (
	"context"
	"time"

	"github.com/goonsorrow/finance-tracker-api/internal/models"
	"github.com/redis/go-redis/v9"
)

type Authorization interface {
	CacheRefreshToken(ctx context.Context, key string, refreshTTL time.Duration) error
	CheckRefreshToken(ctx context.Context, key string) (int, error)
	DeleteRefreshToken(ctx context.Context, key string) error
}

type Category interface {
	GetCatgory(ctx context.Context, id int) (*models.Category, error)
}

type Cache struct {
	Authorization
	Category
}

func NewCache(rdb *redis.Client) *Cache {
	return &Cache{
		Authorization: NewAuthRedis(rdb),
		Category:      NewCategoryRedis(rdb),
	}
}
