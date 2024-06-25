package network

import (
	"context"

	dnetwork "github.com/docker/docker/api/types/network"
)

type CreateOption = dnetwork.CreateOptions

func (n *network) Create(ctx context.Context, name string, opts ...func(*CreateOption)) (dnetwork.CreateResponse, error) {
	opt := &CreateOption{}
	for _, o := range opts {
		o(opt)
	}

	return n.client.NetworkCreate(ctx, name, *opt)
}
