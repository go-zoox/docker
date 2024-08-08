package info

import (
	"context"

	"github.com/docker/docker/api/types/system"
)

func (i *info) Get(ctx context.Context) (system.Info, error) {
	return i.client.Info(ctx)
}
