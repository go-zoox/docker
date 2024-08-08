package info

import (
	"context"

	"github.com/docker/docker/api/types"
)

func (i *info) Disk(ctx context.Context, opts ...func(cfg *types.DiskUsageOptions)) (types.DiskUsage, error) {
	cfg := &types.DiskUsageOptions{}
	for _, o := range opts {
		o(cfg)
	}

	return i.client.DiskUsage(ctx, *cfg)
}
