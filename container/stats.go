package container

import (
	"context"

	"github.com/docker/docker/api/types"
)

type StatsOptions struct {
	Stream bool
}

// Status retrieves the status of a container.
func (c *container) Stats(ctx context.Context, id string, opts ...func(opt *StatsOptions)) (*types.ContainerStats, error) {
	opt := &StatsOptions{}
	for _, o := range opts {
		o(opt)
	}

	stats, err := c.client.ContainerStats(ctx, id, opt.Stream)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}
