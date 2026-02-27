package component

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"
)

type RetryNode struct {
	base[*RetryNode]
}

func (r *RetryNode) Run(ctx context.Context) error {
	r.logger.Info("supervisor", "status", ReadyCheckStart, "component", r.Component.Name())
	err := r.Component.Ready(ctx)
	r.logger.Info("supervisor", "status", ReadyCheckFinish, "component", r.Component.Name())

	if err != nil {
		return err
	}
	g, errCtx := errgroup.WithContext(ctx)
	wg := sync.WaitGroup{}
	wg.Add(1)
	g.Go(func() error {
		r.logger.Info("supervisor", "status", Run, "component", r.Component.Name())
		wg.Done()
		return r.Component.Run(ctx)
	})
	wg.Wait()

	for _, node := range r.Nodes {
		g.Go(func() error {
			return node.Run(errCtx)
		})
	}
	err = g.Wait()
	for _, component := range r.Nodes {
		err := component.Shutdown(ctx)
		if err != nil {
			return err
		}
	}
	r.logger.Info("supervisor", "status", ShutdownStart, "component", r.Component.Name())
	err = r.Component.Shutdown(ctx)
	r.logger.Info("supervisor", "status", ShutdownFinish, "component", r.Component.Name())
	return err
}
