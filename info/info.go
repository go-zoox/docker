package info

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/system"
	"github.com/docker/docker/client"
)

// Info is a client for the info API
type Info interface {
	Get(ctx context.Context) (system.Info, error)
	//
	Version(ctx context.Context) (types.Version, error)
	//
	Disk(ctx context.Context, opts ...func(cfg *types.DiskUsageOptions)) (types.DiskUsage, error)
}

type info struct {
	client *client.Client
}

// New creates a info client
func New(client *client.Client) Info {
	return &info{
		client: client,
	}
}
