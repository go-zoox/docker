package network

import (
	"context"

	"github.com/docker/docker/api/types"
)

type CreateOption = types.NetworkCreate

func (n *network) Create(ctx context.Context, name string, opts ...func(*CreateOption)) (types.NetworkCreateResponse, error) {
	opt := &CreateOption{
		CheckDuplicate: true,
	}
	for _, o := range opts {
		o(opt)
	}

	return n.client.NetworkCreate(ctx, name, *opt)
}
