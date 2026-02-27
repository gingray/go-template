package component

import (
	"context"
	"fmt"
)

type RetryComponent struct {
	BaseComponent
	Retry     int
	Component Component
}

func (r *RetryComponent) Name() string {
	return "retry"
}

func (r *RetryComponent) Run(ctx context.Context) error {
	for i := 0; i < r.Retry; i++ {
		errorCh := make(chan error)
		go func() {
			err := r.Component.Run(ctx)
			if err != nil {
				errorCh <- err
			}
		}()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errorCh:
			fmt.Println(err)
			err = r.Component.Shutdown(ctx)
			fmt.Println(err)
			err = r.Component.Ready(ctx)
			fmt.Println(err)
		}
	}
	return nil
}

func NewRetryComponent(component Component, retry int) *RetryComponent {
	retryComponent := &RetryComponent{Component: component, Retry: retry}
	retryComponent.AddReadyHandler(component.Ready)
	retryComponent.AddShutdownHandler(component.Shutdown)
	return retryComponent
}
