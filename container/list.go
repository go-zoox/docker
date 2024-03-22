package container

import (
	"context"

	tc "github.com/docker/docker/api/types/container"
	"github.com/go-zoox/docker/entity"
)

type ListOptions = tc.ListOptions

func (c *container) List(ctx context.Context, opts ...func(opt *ListOptions)) ([]entity.Container, error) {
	opt := &ListOptions{}
	for _, o := range opts {
		o(opt)
	}

	return c.client.ContainerList(ctx, *opt)
}
