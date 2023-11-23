package volume

import (
	"context"

	vo "github.com/docker/docker/api/types/volume"
)

type InspectOption struct {
}

func (n *volume) Inspect(ctx context.Context, id string, opts ...func(*InspectOption)) (vo.Volume, error) {
	opt := &InspectOption{}
	for _, o := range opts {
		o(opt)
	}

	return n.client.VolumeInspect(ctx, id)
}
