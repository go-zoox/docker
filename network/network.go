package network

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Network interface {
	List(ctx context.Context, opts ...func(*ListOption)) ([]types.NetworkResource, error)
	Inspect(ctx context.Context, id string, opts ...func(*InspectOption)) (types.NetworkResource, error)
	//
	Create(ctx context.Context, name string, opts ...func(*CreateOption)) (types.NetworkCreateResponse, error)
	Remove(ctx context.Context, id string, opts ...func(opt *RemoveOption)) error
	//
	Prune(ctx context.Context, opts ...func(opt *PruneOption)) (types.NetworksPruneReport, error)
}

type network struct {
	client *client.Client
}

// New creates a docker image client
func New(client *client.Client) Network {
	return &network{
		client: client,
	}
}
