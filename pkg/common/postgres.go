package common

import (
	"context"
	"database/sql"
	"fmt"
)

type PostgresConfig struct {
	Host     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	Port     int    `env:"POSTGRES_PORT" envDefault:"5432"`
	User     string `env:"POSTGRES_USER" envDefault:"postgres"`
	Password string `env:"POSTGRES_PASSWORD" envDefault:""`
	DBName   string `env:"POSTGRES_DB" envDefault:"postgres"`
}

func (p PostgresConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", p.User, p.Password, p.Host, p.Port, p.DBName)
}

func (a *App) WithPostgres(cfg *PostgresConfig) error {
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
