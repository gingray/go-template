package component

import (
	"context"
	"fmt"
)

type ShutdownComponent struct {
	BaseComponent
	shutdownCh chan struct{}
}

func NewShutdownComponent(shutdownCh chan struct{}) *ShutdownComponent {
	return &ShutdownComponent{shutdownCh: shutdownCh}
}

func (s *ShutdownComponent) Name() string {
	return "shutdown"
}

func (s *ShutdownComponent) ShutdownContext(ctx context.Context) context.Context {
	shutDownCtx, cancel := context.WithCancel(ctx)
	go func() {
		<-s.shutdownCh
		cancel()
	}()
	return shutDownCtx
}

func (s *ShutdownComponent) Run(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-s.shutdownCh:
		return fmt.Errorf("shutdown")
	}
}
