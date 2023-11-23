package image

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

type PruneOption = filters.Args

func (i *image) Prune(ctx context.Context, opts ...func(opt *PruneOption)) (types.ImagesPruneReport, error) {
	opt := &PruneOption{}
	for _, o := range opts {
		o(opt)
	}

	return i.client.ImagesPrune(ctx, *opt)
}
