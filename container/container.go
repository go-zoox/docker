package container

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	dc "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/go-zoox/docker/entity"
)

// Container is the docker container client interface
type Container interface {
	List(ctx context.Context, opts ...func(opt *ListOptions)) ([]entity.Container, error)
	//
	Create(ctx context.Context, opts ...func(opt *CreateOptions)) (dc.CreateResponse, error)
	Update(ctx context.Context, id string, opts ...func(opt *UpdateOptions)) (dc.ContainerUpdateOKBody, error)
	Remove(ctx context.Context, id string, opts ...func(opt *RemoveOptions)) error
	//
	Start(ctx context.Context, id string, opts ...func(opt *StartOptions)) error
	Stop(ctx context.Context, id string, opts ...func(opt *StopOptions)) error
	Restart(ctx context.Context, id string, opts ...func(opt *RestartOptions)) error
	//
	Inspect(ctx context.Context, id string, opts ...func(opt *InspectOptions)) (*types.ContainerJSON, error)
	//
	Stats(ctx context.Context, id string, opts ...func(opt *StatsOptions)) (*dc.StatsResponseReader, error)
	//
	Logs(ctx context.Context, id string, opts ...func(opt *LogsConfig)) (io.ReadCloser, error)
	//
	Exec(ctx context.Context, id string, opts ...func(opt *ExecOptions)) (*ExecTerm, error)
	//
	Run(ctx context.Context, opts ...func(opt *RunOptions)) error
}

type container struct {
	client *client.Client
}

// New creates a docker container client
func New(client *client.Client) Container {
	return &container{
		client: client,
	}
}

// cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
// if err != nil {
// 	return err
// }
