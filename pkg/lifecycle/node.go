package lifecycle

import (
	"context"
	"errors"

	"github.com/gingray/go-template/pkg/config"
)

const (
	ReadyCheckStart  = "ready-check-start"
	ReadyCheckFinish = "ready-check-finish"
	Run              = "running"
	ShutdownStart    = "shutdown-start"
	ShutdownFinish   = "shutdown-finish"
)

var ComponentStopErr = errors.New("component stop")

type strategy interface {
	Process(ctx context.Context, baseNode *Node) error
}

func GetNodeCreator(logger config.Logger) func(component Component) *Node {
	var creator func(component Component) *Node
	creator = func(component Component) *Node {
		return &Node{Component: component, Nodes: []*Node{}, logger: logger, strategy: &DefaultStrategy{},
			nodeCreator: creator}
	}
	return creator
}

func DefaultRoot(logger config.Logger) *Node {
	creator := GetNodeCreator(logger)
	root := creator(NewRootComponent())
	return root
}

func (n *Node) Then(component ...Component) {
	for _, item := range component {
		switch c := item.(type) {
		case *Node:
			n.Nodes = append(n.Nodes, c)
		default:
			n.Nodes = append(n.Nodes, n.nodeCreator(c))
		}
	}
}

func (n *Node) ThenLast(component ...Component) *Node {
	n.Then(component...)
	return n.Nodes[len(n.Nodes)-1]
}

func (n *Node) ThenFirst(component ...Component) *Node {
	n.Then(component...)
	return n.Nodes[0]
}

func (n *Node) ThenNth(nth int, component ...Component) *Node {
	n.Then(component...)
	return n.Nodes[nth-1]
}

func (n *Node) asNode(component Component) *Node {
	return n.nodeCreator(component)
}

func (n *Node) Shutdown(ctx context.Context) error {
	for _, node := range n.Nodes {
		err := node.Shutdown(ctx)
		if err != nil {
			return err
		}
	}
	return n.Component.Shutdown(ctx)
}

type Node struct {
	Component   Component
	Nodes       []*Node
	logger      config.Logger
	strategy    strategy
	nodeCreator func(component Component) *Node
}

func (n *Node) Name() string {
	return n.Component.Name()
}

func (n *Node) Ready(ctx context.Context) error {
	return n.Component.Ready(ctx)
}

func (n *Node) Run(ctx context.Context) error {
	n.logger.Info("supervisor", "status", ReadyCheckStart, "component", n.Component.Name())
	err := n.Component.Ready(ctx)
	n.logger.Info("supervisor", "status", ReadyCheckFinish, "component", n.Component.Name())

	if err != nil {
		return err
	}
	return n.strategy.Process(ctx, n)
}
