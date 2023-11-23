package container

import (
	"context"

	"github.com/docker/docker/api/types"
)

type RemoveOptions = types.ContainerRemoveOptions

// Remove removes a container.
func (c *container) Remove(ctx context.Context, id string, opts ...func(opt *RemoveOptions)) error {
	opt := &RemoveOptions{
		Force: true,
	}
	for _, o := range opts {
		o(opt)
	}

	return c.client.ContainerRemove(ctx, id, *opt)
}
