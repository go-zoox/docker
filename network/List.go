package network

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

type ListOption struct {
	Filters filters.Args
}

func (n *network) List(ctx context.Context, opts ...func(*ListOption)) ([]types.NetworkResource, error) {
	opt := &ListOption{}
	for _, o := range opts {
		o(opt)
	}

	return n.client.NetworkList(ctx, types.NetworkListOptions{
		Filters: opt.Filters,
	})
}
