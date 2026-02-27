package component

import (
	"context"
	"log/slog"
	"sync"

	"golang.org/x/sync/errgroup"
)

const (
	ReadyCheckStart  = "ready-check-start"
	ReadyCheckFinish = "ready-check-finish"
	Run              = "running"
	ShutdownStart    = "shutdown-start"
	ShutdownFinish   = "shutdown-finish"
)

type Node interface {
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type Supervisor struct {
	logger *slog.Logger
}

func NewSupervisor(logger *slog.Logger) *Supervisor {
	return &Supervisor{logger: logger}
}

type base[T Node] struct {
	Component Component
	Nodes     []Node
	logger    *slog.Logger
}

func (b *base[T]) AddNode(node Node) {
	b.Nodes = append(b.Nodes, node)
}

func (b *base[T]) Shutdown(ctx context.Context) error {
	for _, node := range b.Nodes {
		err := node.Shutdown(ctx)
		if err != nil {
			return err
		}
	}
	return b.Component.Shutdown(ctx)
}

type BaseNode struct {
	base[*BaseNode]
}

func (s *Supervisor) CreateRootNode() *BaseNode {
	return &BaseNode{base[*BaseNode]{Component: NewRootComponent(), Nodes: []Node{}, logger: s.logger}}
}

func (s *Supervisor) CreateNode(component Component) *BaseNode {
	return &BaseNode{base[*BaseNode]{Component: component, Nodes: []Node{}, logger: s.logger}}
}

func (n *BaseNode) Run(ctx context.Context) error {
	n.logger.Info("supervisor", "status", ReadyCheckStart, "component", n.Component.Name())
	err := n.Component.Ready(ctx)
	n.logger.Info("supervisor", "status", ReadyCheckFinish, "component", n.Component.Name())

	if err != nil {
		return err
	}
	g, errCtx := errgroup.WithContext(ctx)
	wg := sync.WaitGroup{}
	wg.Add(1)
	g.Go(func() error {
		n.logger.Info("supervisor", "status", Run, "component", n.Component.Name())
		wg.Done()
		return n.Component.Run(ctx)
	})
	wg.Wait()

	for _, node := range n.Nodes {
		g.Go(func() error {
			return node.Run(errCtx)
		})
	}
	err = g.Wait()
	for _, component := range n.Nodes {
		err := component.Shutdown(ctx)
		if err != nil {
			return err
		}
	}
	n.logger.Info("supervisor", "status", ShutdownStart, "component", n.Component.Name())
	err = n.Component.Shutdown(ctx)
	n.logger.Info("supervisor", "status", ShutdownFinish, "component", n.Component.Name())
	return err
}
