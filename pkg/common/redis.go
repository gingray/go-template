package common

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host     string `env:"REDIS_HOST" envDefault:"localhost"`
	Password string `env:"REDIS_PASSWORD" envDefault:""`
	Port     int    `env:"REDIS_PORT" envDefault:"6379"`
	Db       int    `env:"REDIS_DB" envDefault:"0"`
}

func (a *App) WithRedis(cfg *RedisConfig) error {
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

func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}
