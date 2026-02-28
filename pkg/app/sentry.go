package app

import (
	"github.com/getsentry/sentry-go"
	"github.com/gingray/go-template/pkg/config"
)

func InitSentry(cfg *config.SentryConfig) error {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         cfg.DSN,
		Environment: cfg.Env,
		Release:     cfg.Release,
	})
	return err
}
