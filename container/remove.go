package container

import (
	"context"

	tc "github.com/docker/docker/api/types/container"
)

type RemoveOptions = tc.RemoveOptions

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
