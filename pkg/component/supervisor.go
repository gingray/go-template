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

type Node struct {
	Component Component
	Nodes     []*Node
	logger    *slog.Logger
}

func NewSupervisor(logger *slog.Logger) *Node {
	return &Node{Component: NewRootComponent(), Nodes: []*Node{}, logger: logger}
}

func (n *Node) NewNode(component Component) *Node {
	return &Node{
		Component: component,
		Nodes:     []*Node{},
		logger:    n.logger,
	}
}

func (n *Node) AddComponent(component Component) *Node {
	newNode := n.NewNode(component)
	n.Nodes = append(n.Nodes, newNode)
	return newNode
}
func (n *Node) Run(ctx context.Context) error {
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
		err := component.shutdown(ctx)
		if err != nil {
			return err
		}
	}
	n.logger.Info("supervisor", "status", ShutdownStart, "component", n.Component.Name())
	err = n.Component.Shutdown(ctx)
	n.logger.Info("supervisor", "status", ShutdownFinish, "component", n.Component.Name())
	return err
}

func (n *Node) shutdown(ctx context.Context) error {
	for _, node := range n.Nodes {
		err := node.shutdown(ctx)
		if err != nil {
			return err
		}
	}
	return n.Component.Shutdown(ctx)
}
