package image

import (
	"context"

	ti "github.com/docker/docker/api/types/image"
)

type RemoveConfig struct {
	Force         bool
	PruneChildren bool
}

func (i *image) Remove(ctx context.Context, id string, opts ...func(opt *RemoveConfig)) ([]ti.DeleteResponse, error) {
	opt := &RemoveConfig{}
	for _, o := range opts {
		o(opt)
	}

	return i.client.ImageRemove(ctx, id, ti.RemoveOptions{
		Force:         opt.Force,
		PruneChildren: opt.PruneChildren,
	})
}
