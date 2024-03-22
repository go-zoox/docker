package image

import (
	"context"

	"github.com/docker/docker/api/types"
	ti "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/go-zoox/docker/entity"
)

// Image is the docker image client interface
type Image interface {
	List(ctx context.Context, opts ...func(opt *ListOption)) (images []entity.Image, err error)
	Inspect(ctx context.Context, id string, opts ...func(opt *InspectOption)) (*types.ImageInspect, error)
	Remove(ctx context.Context, id string, opts ...func(opt *RemoveOption)) ([]ti.DeleteResponse, error)
	//
	Build(ctx context.Context, src string, opts ...func(opt *BuildOption)) error
	Pull(ctx context.Context, name string, opts ...func(opt *PullOption)) error
	// Push(ctx context.Context, opts ...func(opt *PushOption)) error
	Tag(ctx context.Context, source, target string) error
	//
	Prune(ctx context.Context, opts ...func(opt *PruneOption)) (types.ImagesPruneReport, error)
}

type image struct {
	client *client.Client
}

// New creates a docker image client
func New(client *client.Client) Image {
	return &image{
		client: client,
	}
}
