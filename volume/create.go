package volume

import (
	"context"

	vo "github.com/docker/docker/api/types/volume"
)

type CreateOption = vo.CreateOptions

func (n *volume) Create(ctx context.Context, opts ...func(*CreateOption)) (vo.Volume, error) {
	cfg := &CreateOption{}
	for _, o := range opts {
		o(cfg)
	}

	return n.client.VolumeCreate(ctx, *cfg)
}
