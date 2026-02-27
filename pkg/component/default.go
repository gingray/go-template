package component

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

type DefaultStrategy struct {
}

func (i *DefaultStrategy) Process(ctx context.Context, n *Node) error {
	g, errCtx := errgroup.WithContext(ctx)
	wg := sync.WaitGroup{}
	wg.Add(1)
	g.Go(func() error {
		n.logger.Info("supervisor", "status", Run, "component", n.Component.Name())
		wg.Done()
		return n.Component.Run(errCtx)
	})
	wg.Wait()

	for _, node := range n.Nodes {
		g.Go(func() error {
			return node.Run(errCtx)
		})
	}
	err := g.Wait()
	for _, component := range n.Nodes {
		compErr := component.Shutdown(ctx)
		err = errors.Join(err, compErr)
	}
	n.logger.Info("supervisor", "status", ShutdownStart, "component", n.Component.Name())
	stopErr := fmt.Errorf("component: %s, %w", n.Component.Name(), componentStopErr)
	err = errors.Join(err, n.Component.Shutdown(ctx), stopErr)
	n.logger.Info("supervisor", "status", ShutdownFinish, "component", n.Component.Name())
	return err
}
