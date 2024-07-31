package image

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

type PruneConfig = filters.Args

func (i *image) Prune(ctx context.Context, opts ...func(opt *PruneConfig)) (types.ImagesPruneReport, error) {
	opt := &PruneConfig{}
	for _, o := range opts {
		o(opt)
	}

	return i.client.ImagesPrune(ctx, *opt)
}
