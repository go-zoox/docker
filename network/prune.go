package network

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

type PruneOption struct {
	Filters filters.Args
}

func (n *network) Prune(ctx context.Context, opts ...func(opt *PruneOption)) (types.NetworksPruneReport, error) {
	opt := &PruneOption{}
	for _, o := range opts {
		o(opt)
	}

	return n.client.NetworksPrune(ctx, opt.Filters)
}
