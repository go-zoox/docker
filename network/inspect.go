package network

import (
	"context"

	dn "github.com/docker/docker/api/types/network"
	"github.com/go-zoox/docker/entity"
)

type InspectOption struct {
	Scope   string
	Verbose bool
}

func (n *network) Inspect(ctx context.Context, id string, opts ...func(*InspectOption)) (entity.Network, error) {
	opt := &InspectOption{}
	for _, o := range opts {
		o(opt)
	}

	return n.client.NetworkInspect(ctx, id, dn.InspectOptions{})
}
