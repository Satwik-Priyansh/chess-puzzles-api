package db

import (
	"chess-puzzles-api/config"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(cfg *config.EnvConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		slog.Error("Failed to parse database configuration string", "error", err)
		return nil, err
	}
	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.MaxConnLifetime = 1 * time.Hour

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		slog.Error("Failed to initialize database connection pool", "error", err)
		return nil, err
	}
	err = pool.Ping(ctx)
	if err != nil {
		slog.Error("Database connection ping failed", "error", err)
		return nil, err
	}
	return pool, nil
}
