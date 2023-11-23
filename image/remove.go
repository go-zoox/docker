package image

import (
	"context"

	"github.com/docker/docker/api/types"
)

type RemoveOption struct {
	Force         bool
	PruneChildren bool
}

func (i *image) Remove(ctx context.Context, id string, opts ...func(opt *RemoveOption)) ([]types.ImageDeleteResponseItem, error) {
	opt := &RemoveOption{}
	for _, o := range opts {
		o(opt)
	}

	return i.client.ImageRemove(ctx, id, types.ImageRemoveOptions{
		Force:         opt.Force,
		PruneChildren: opt.PruneChildren,
	})
}
