package network

import (
	"context"
)

type RemoveOption struct {
}

func (n *network) Remove(ctx context.Context, id string, opts ...func(opt *RemoveOption)) error {
	opt := &RemoveOption{}
	for _, o := range opts {
		o(opt)
	}

	return n.client.NetworkRemove(ctx, id)
}
