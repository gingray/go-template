package common

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/gingray/go-template/pkg/component"
	"github.com/redis/go-redis/v9"
	"github.com/twmb/franz-go/pkg/kgo"
)

type App struct {
	component.BaseComponent
	PGdb       *sql.DB
	Rdb        *redis.Client
	Kafka      *kgo.Client
	HttpRouter *gin.Engine
	Logger     *slog.Logger
}

func (a *App) Name() string {
	return "app"
}

func NewApp(cfg *Config) (*App, error) {
	app := &App{}
	err := app.WithLogger()
	if err != nil {
		return nil, err
	}

	err = app.WithHTTPRouter()
	if err != nil {
		return nil, err
	}

	err = app.WithPostgres(&cfg.PostgresConfig)
	if err != nil {
		return nil, err
	}
	//err = app.WithRedis(&cfg.RedisConfig)
	//if err != nil {
	//	return nil, err
	//}
	//
	//err = app.WithKafka(&cfg.KafkaConfig)
	//if err != nil {
	//	return nil, err
	//}
	//err = InitSentry(&cfg.SentryConfig)
	//if err != nil {
	//	return nil, err
	//}
	//err = InitDataDog(&cfg.DataDogConfig)
	//if err != nil {
	//	return nil, err
	//}
	return app, nil
}

func (a *App) Run(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	}
}
