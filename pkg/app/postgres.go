package app

import (
	"context"
	"fmt"

	"github.com/gingray/go-template/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (a *App) WithPostgres(cfg *config.PostgresConfig) error {
	pgPool, err := pgxpool.New(context.Background(), cfg.DSN())
	a.PgPool = pgPool
	a.AddReadyHandler(func(ctx context.Context) error {
		err := pgPool.Ping(ctx)
		if err != nil {
			err = fmt.Errorf("ping postgres: %w", err)
		}
		return err
	})
	a.AddShutdownHandler(func(ctx context.Context) error {
		pgPool.Close()
		return nil
	})
	return err
}
