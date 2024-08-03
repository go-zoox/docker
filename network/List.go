package network

import (
	"context"

	"github.com/docker/docker/api/types/filters"
	dn "github.com/docker/docker/api/types/network"
	"github.com/go-zoox/docker/entity"
)

type ListOption struct {
	Filters filters.Args
}

func (n *network) List(ctx context.Context, opts ...func(*ListOption)) ([]entity.Network, error) {
	opt := &ListOption{}
	for _, o := range opts {
		o(opt)
	}

	return n.client.NetworkList(ctx, dn.ListOptions{
		Filters: opt.Filters,
	})
}
