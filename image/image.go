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
	List(ctx context.Context, opts ...func(cfg *ListConfig)) (images []entity.Image, err error)
	Inspect(ctx context.Context, id string, opts ...func(cfg *InspectConfig)) (*types.ImageInspect, error)
	Remove(ctx context.Context, id string, opts ...func(cfg *RemoveConfig)) ([]ti.DeleteResponse, error)
	//
	History(ctx context.Context, id string) ([]entity.ImageHistory, error)
	//
	Build(ctx context.Context, src string, opts ...func(cfg *BuildConfig)) error
	//
	Pull(ctx context.Context, name string, opts ...func(cfg *PullConfig)) error
	Push(ctx context.Context, name string, opts ...func(cfg *PushConfig)) error
	//
	Tag(ctx context.Context, source, target string) error
	//
	Prune(ctx context.Context, opts ...func(cfg *PruneConfig)) (types.ImagesPruneReport, error)
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
