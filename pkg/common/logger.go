package common

import (
	"log/slog"
	"os"
)

func (a *App) WithLogger() error {
	a.Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return nil
}
