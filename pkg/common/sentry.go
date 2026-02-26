package common

import "github.com/getsentry/sentry-go"

type SentryConfig struct {
	DSN     string `env:"SENTRY_DSN" envDefault:""`
	Env     string `env:"SENTRY_ENV" envDefault:"dev"`
	Release string `env:"SENTRY_RELEASE" envDefault:"my-app@1.0.0"`
}

func InitSentry(cfg *SentryConfig) error {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         cfg.DSN,
		Environment: cfg.Env,
		Release:     cfg.Release,
	})
	return err
}
