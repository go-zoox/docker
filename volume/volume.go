package volume

import (
	"context"

	"github.com/docker/docker/api/types"
	vo "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/go-zoox/docker/entity"
)

type Volume interface {
	List(ctx context.Context, opts ...func(*ListOption)) ([]entity.Volume, error)
	Inspect(ctx context.Context, id string, opts ...func(*InspectOption)) (vo.Volume, error)
	//
	Create(ctx context.Context, opts ...func(*CreateOption)) (vo.Volume, error)
	Remove(ctx context.Context, id string, opts ...func(opt *RemoveOption)) error
	//
	Prune(ctx context.Context, opts ...func(opt *PruneOption)) (types.VolumesPruneReport, error)
}

type volume struct {
	client *client.Client
}

// New creates a docker image client
func New(client *client.Client) Volume {
	return &volume{
		client: client,
	}
}
