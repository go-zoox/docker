package container

import (
	"context"

	co "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

type CreateOptions struct {
	Name string
	//
	Container *co.Config
	Host      *co.HostConfig
	Network   *network.NetworkingConfig
	Platform  *specs.Platform
}

// Create creates a container
func (c *container) Create(ctx context.Context, opts ...func(opt *CreateOptions)) (co.CreateResponse, error) {
	optsX := &CreateOptions{
		Container: &co.Config{},
		Host:      &co.HostConfig{},
		Network:   &network.NetworkingConfig{},
		Platform:  &specs.Platform{},
	}
	for _, opt := range opts {
		opt(optsX)
	}

	return c.client.ContainerCreate(ctx, optsX.Container, optsX.Host, optsX.Network, optsX.Platform, optsX.Name)
}
