package network

import (
	"context"

	"github.com/docker/docker/api/types"
)

type InspectOption struct {
	Scope   string
	Verbose bool
}

func (n *network) Inspect(ctx context.Context, id string, opts ...func(*InspectOption)) (types.NetworkResource, error) {
	opt := &InspectOption{}
	for _, o := range opts {
		o(opt)
	}

	return n.client.NetworkInspect(ctx, id, types.NetworkInspectOptions{})
}
