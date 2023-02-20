package container

import (
	"context"

	dc "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// Stop stops a container
func Stop(id string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	timeout := 60
	if err := cli.ContainerStop(ctx, id, dc.StopOptions{
		Timeout: &timeout,
	}); err != nil {
		return err
	}

	return nil
}
