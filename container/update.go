package container

import (
	"context"

	co "github.com/docker/docker/api/types/container"
)

type UpdateOptions struct {
	co.Resources
	RestartPolicy co.RestartPolicy
}

// Update updates a container.
func (c *container) Update(ctx context.Context, id string, opts ...func(opt *UpdateOptions)) (co.ContainerUpdateOKBody, error) {
	opt := &UpdateOptions{}
	for _, o := range opts {
		o(opt)
	}

	return c.client.ContainerUpdate(ctx, id, co.UpdateConfig{
		RestartPolicy: opt.RestartPolicy,
		Resources:     opt.Resources,
	})
}
