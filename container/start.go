package container

import (
	"context"

	"github.com/docker/docker/api/types"
)

type StartOptions = types.ContainerStartOptions

// Start starts a container
func (c *container) Start(ctx context.Context, id string, opts ...func(opt *StartOptions)) error {
	optsX := &StartOptions{}
	for _, opt := range opts {
		opt(optsX)
	}

	return c.client.ContainerStart(ctx, id, *optsX)
}
