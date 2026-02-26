package common

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host     string `env:"REDIS_HOST" envDefault:"localhost"`
	Password string `env:"REDIS_PASSWORD" envDefault:""`
	Port     int    `env:"REDIS_PORT" envDefault:"6379"`
	Db       int    `env:"REDIS_DB" envDefault:"0"`
}

func NewRedis(cfg *RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.Db,
	})
	return rdb, nil
}

func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}
