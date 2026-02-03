package cache

import (
	"context"

	"github.com/goonsorrow/finance-tracker-api/internal/models"
	"github.com/redis/go-redis/v9"
)

type CategoryRedis struct {
	rdb *redis.Client
}

func NewCategoryRedis(rdb *redis.Client) *CategoryRedis {
	return &CategoryRedis{rdb: rdb}
}

func (c CategoryRedis) GetCatgory(ctx context.Context, id int) (*models.Category, error) {
	return &models.Category{}, nil
}

//DelCategory

//SetCategory

//UpdateCategory
