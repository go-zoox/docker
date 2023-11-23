package container

import (
	"context"

	co "github.com/docker/docker/api/types/container"
)

type RestartOptions = StopOptions

// Restart restarts a container
func (c *container) Restart(ctx context.Context, id string, opts ...func(opt *RestartOptions)) error {
	optsX := &RestartOptions{
		Timeout: 10,
	}
	for _, opt := range opts {
		opt(optsX)
	}

	return c.client.ContainerRestart(ctx, id, co.StopOptions{
		Timeout: &optsX.Timeout,
		Signal:  optsX.Signal,
	})
}
