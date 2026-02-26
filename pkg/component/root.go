package component

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

type RootComponent struct {
	BaseComponent
}

func (r *RootComponent) ComponentName() string {
	return "root"
}

func (r *RootComponent) Run(ctx context.Context) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	return nil
}
