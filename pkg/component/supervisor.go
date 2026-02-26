package component

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Node struct {
	Component Component
	Nodes     []*Node
}

func NewNode(component Component) *Node {
	return &Node{
		Component: component,
		Nodes:     []*Node{},
	}
}

func (n *Node) AddComponent(component Component) *Node {
	newNode := NewNode(component)
	n.Nodes = append(n.Nodes, newNode)
	return newNode
}
func (n *Node) Run(ctx context.Context) error {
	err := n.Component.Ready(ctx)
	if err != nil {
		return err
	}
	g, errCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return n.Component.Run(ctx)
	})

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
	err = n.Component.Shutdown(ctx)

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
