package app

import (
	"context"

	"github.com/gingray/go-template/pkg/config"
	"github.com/redis/go-redis/v9"
)

func (a *App) WithRedis(cfg *config.RedisConfig) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.Db,
	})
	a.Rdb = rdb
	a.AddReadyHandler(func(ctx context.Context) error {
		return rdb.Ping(ctx).Err()
	})
	a.AddShutdownHandler(func(ctx context.Context) error {
		return rdb.Close()
	})

	return nil
}
