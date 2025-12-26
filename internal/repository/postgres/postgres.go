package postgres

import (
	"context"
	"fmt"

	"github.com/adexcell/go-tutorial/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(ctx context.Context, cfg config.PostgresConfig) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать poolConfig: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.MaxOpenConns)
	poolConfig.MaxConnLifetime = cfg.ConnMaxLifetime
	poolConfig.MinIdleConns = int32(cfg.MaxIdleConns)
	poolConfig.MaxConnIdleTime = cfg.ConnMaxIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("не пингуется, %w", err)
	}

	return pool, nil
}
