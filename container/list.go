package container

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/go-zoox/docker/entity"
)

type ListOptions = types.ContainerListOptions

func (c *container) List(ctx context.Context, opts ...func(opt *ListOptions)) ([]entity.Container, error) {
	opt := &ListOptions{}
	for _, o := range opts {
		o(opt)
	}

	return c.client.ContainerList(ctx, *opt)
}
