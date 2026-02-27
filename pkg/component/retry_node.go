package component

import "context"

type RetryNode struct {
	BaseNode
}

func (r *RetryNode) Run(ctx context.Context) error {
	return nil
}
