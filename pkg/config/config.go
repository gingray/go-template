package config

import (
	"fmt"

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

type PostgresConfig struct {
	Host     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	Port     int    `env:"POSTGRES_PORT" envDefault:"5432"`
	User     string `env:"POSTGRES_USER" envDefault:"postgres"`
	Password string `env:"POSTGRES_PASSWORD" envDefault:""`
	DBName   string `env:"POSTGRES_DB" envDefault:"postgres"`
}

type RedisConfig struct {
	Host     string `env:"REDIS_HOST" envDefault:"localhost"`
	Password string `env:"REDIS_PASSWORD" envDefault:""`
	Port     int    `env:"REDIS_PORT" envDefault:"6379"`
	Db       int    `env:"REDIS_DB" envDefault:"0"`
}

type KafkaConfig struct {
	Brokers       []string `env:"KAFKA_BROKERS" envDefault:"localhost:9092"`
	ConsumerGroup string   `env:"KAFKA_CONSUMER_GROUP" envDefault:"test"`
	ConsumeTopics []string `env:"KAFKA_CONSUME_TOPICS" envDefault:"test"`
}

type HTTPServiceConfig struct {
	Port int `env:"HTTP_PORT" envDefault:"3000"`
}

type SentryConfig struct {
	DSN     string `env:"SENTRY_DSN" envDefault:""`
	Env     string `env:"SENTRY_ENV" envDefault:"dev"`
	Release string `env:"SENTRY_RELEASE" envDefault:"my-app@1.0.0"`
}

type DataDogConfig struct {
	AgentAddr        string `env:"DD_AGENT_ADDR"       envDefault:"localhost:8126"`
	ServiceName      string `env:"DD_SERVICE"          envDefault:"my-service"`
	Environment      string `env:"DD_ENV"              envDefault:"dev"`
	Version          string `env:"DD_VERSION"          envDefault:"1.0.0"`
	Enabled          bool   `env:"DD_ENABLED"          envDefault:"false"`
	ProfilingEnabled bool   `env:"DD_PROFILING_ENABLED" envDefault:"false"`
}

func (p PostgresConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", p.User, p.Password, p.Host, p.Port, p.DBName)
}

func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
