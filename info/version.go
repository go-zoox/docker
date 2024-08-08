package info

import (
	"context"

	"github.com/docker/docker/api/types"
)

func (i *info) Version(ctx context.Context) (types.Version, error) {
	return i.client.ServerVersion(ctx)
}
