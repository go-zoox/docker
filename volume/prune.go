package volume

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

type PruneOption struct {
	Filters filters.Args
}

func (n *volume) Prune(ctx context.Context, opts ...func(opt *PruneOption)) (types.VolumesPruneReport, error) {
	opt := &PruneOption{}
	for _, o := range opts {
		o(opt)
	}

	return n.client.VolumesPrune(ctx, opt.Filters)
}
