package container

import (
	"context"

	dc "github.com/docker/docker/api/types/container"
)

type StatsOptions struct {
	Stream bool
}

// Status retrieves the status of a container.
func (c *container) Stats(ctx context.Context, id string, opts ...func(opt *StatsOptions)) (*dc.StatsResponseReader, error) {
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
