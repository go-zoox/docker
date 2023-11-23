package container

import (
	"context"

	"github.com/go-zoox/docker/entity"
)

type InspectOptions struct {
}

// Inspect inspects a container.
func (c *container) Inspect(ctx context.Context, id string, opts ...func(opt *InspectOptions)) (*entity.ContainerState, error) {
	opt := &InspectOptions{}
	for _, o := range opts {
		o(opt)
	}

	info, err := c.client.ContainerInspect(ctx, id)
	if err != nil {
		return nil, err
	}

	return &info, nil
}
