package cache

import (
	"context"
	"fmt"

	"github.com/adexcell/go-tutorial/internal/config"
	"github.com/redis/go-redis/v9"
)

var Cache *redis.Client

func New(ctx context.Context, cfg config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		MinIdleConns: cfg.MinIdleConns,
		PoolSize:     cfg.PoolSize,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("не удалось запустить redis: %w", err)
	}

	return rdb, nil
}
