package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gingray/go-template/pkg/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func (a *App) WithPostgres(cfg *config.PostgresConfig) error {
	db, err := sql.Open("pgx", cfg.DSN())
	a.PGdb = db
	a.AddReadyHandler(func(ctx context.Context) error {
		err := db.PingContext(ctx)
		if err != nil {
			err = fmt.Errorf("ping postgres: %w", err)
		}
		return err
	})
	a.AddShutdownHandler(func(ctx context.Context) error {
		err := db.Close()
		if err != nil {
			err = fmt.Errorf("shutdown postgres: %w", err)
		}
		return err
	})
	return err
}
