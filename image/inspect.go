package image

import (
	"context"

	"github.com/docker/docker/api/types"
)

type InspectOption struct {
}

func (i *image) Inspect(ctx context.Context, id string, opts ...func(opt *InspectOption)) (*types.ImageInspect, error) {
	opt := &InspectOption{}
	for _, o := range opts {
		o(opt)
	}

	inspect, _, err := i.client.ImageInspectWithRaw(ctx, id)
	if err != nil {
		return nil, err
	}

	return &inspect, nil
}
