package common

import (
	"database/sql"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/twmb/franz-go/pkg/kgo"
)

type App struct {
	PGdb       *sql.DB
	Rdb        *redis.Client
	Kafka      *kgo.Client
	HttpRouter *gin.Engine
	Logger     *slog.Logger
}

func NewApp(cfg *Config) (*App, error) {
	db, err := NewPostgres(&cfg.PostgresConfig)
	if err != nil {
		return nil, err
	}
	rdb, err := NewRedis(&cfg.RedisConfig)
	if err != nil {
		return nil, err
	}

	kafka, err := NewKafka(&cfg.KafkaConfig)
	if err != nil {
		return nil, err
	}
	err = InitSentry(&cfg.SentryConfig)
	if err != nil {
		return nil, err
	}
	err = InitDataDog(&cfg.DataDogConfig)
	if err != nil {
		return nil, err
	}
	router := NewHTTPRouter()
	logger := NewLogger()
	return &App{PGdb: db, Rdb: rdb, Kafka: kafka, HttpRouter: router, Logger: logger}, nil
}

func (app *App) Start() error {
	return nil
}
