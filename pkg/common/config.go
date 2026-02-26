package common

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	HTTPServiceConfig
	RedisConfig
	PostgresConfig
	KafkaConfig
	SentryConfig
	DataDogConfig
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
