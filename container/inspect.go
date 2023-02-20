package container

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// Inspect inspects a container
func Inspect(id string) (*types.ContainerJSON, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	info, err := cli.ContainerInspect(ctx, id)
	if err != nil {
		return nil, err
	}

	return &info, nil
}
