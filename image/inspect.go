package image

import (
	"context"

	"github.com/docker/docker/api/types"
)

type InspectConfig struct {
}

func (i *image) Inspect(ctx context.Context, id string, opts ...func(opt *InspectConfig)) (*types.ImageInspect, error) {
	cfg := &InspectConfig{}
	for _, o := range opts {
		o(cfg)
	}

	inspect, _, err := i.client.ImageInspectWithRaw(ctx, id)
	if err != nil {
		return nil, err
	}

	return &inspect, nil
}
