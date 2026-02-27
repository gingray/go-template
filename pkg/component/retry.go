package component

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"
)

type RetryStrategy struct {
	Retry int
}

func (r *RetryStrategy) Process(ctx context.Context, n *Node) error {
	g, errCtx := errgroup.WithContext(ctx)
	wg := sync.WaitGroup{}
	wg.Add(1)
	g.Go(func() error {
		n.logger.Info("supervisor", "status", Run, "component", n.Component.Name())
		wg.Done()
		return n.Component.Run(ctx)
	})
	wg.Wait()
	wg = sync.WaitGroup{}
	for _, node := range n.Nodes {
		wg.Add(1)
		go r.runner(errCtx, node, &wg)
	}
	wg.Wait()
	for _, component := range n.Nodes {
		err := component.Shutdown(ctx)
		if err != nil {
			return err
		}
	}
	n.logger.Info("supervisor", "status", ShutdownStart, "component", n.Component.Name())
	err := n.Component.Shutdown(ctx)
	n.logger.Info("supervisor", "status", ShutdownFinish, "component", n.Component.Name())
	return err
}

func (r *RetryStrategy) runner(ctx context.Context, n *Node, wg *sync.WaitGroup) {
	currentTry := 0
	defer wg.Done()
	for i := 0; i < r.Retry; i++ {
		errCh := make(chan error)
		go func() {
			err := n.Run(ctx)
			errCh <- err
		}()
		select {
		case <-ctx.Done():
			return
		case err := <-errCh:
			n.logger.Error("supervisor", "status", ShutdownStart, "component", n.Component.Name(), "error", err, "currentTry", currentTry)
			err = n.Shutdown(ctx)
			n.logger.Error("supervisor", "status", ShutdownFinish, "component", n.Component.Name(), "error", err, "currentTry", currentTry)
			break
		}
		currentTry++
	}
}
