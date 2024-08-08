package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/system"
	"github.com/docker/docker/client"
	"github.com/go-zoox/docker/container"
	"github.com/go-zoox/docker/image"
	"github.com/go-zoox/docker/network"
	"github.com/go-zoox/docker/volume"
)

// Docker is the client of docker
type Docker interface {
	Info(ctx context.Context) (system.Info, error)
	Version(ctx context.Context) (types.Version, error)

	Container() container.Container
	Image() image.Image
	Network() network.Network
	Volume() volume.Volume
}

type docker struct {
	client *client.Client
}

// Config ...
type Config struct {
	// Docker Host
	Server string `json:"server"`
}

// Option ...
type Option func(cfg *Config)

// New creates a docker client.
func New(opts ...Option) (d Docker, err error) {
	cfg := &Config{}
	for _, o := range opts {
		o(cfg)
	}

	var c *client.Client
	if cfg.Server == "" {
		c, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			return nil, err
		}
	} else {
		c, err = client.NewClientWithOpts(client.WithHost(cfg.Server), client.WithAPIVersionNegotiation())
		if err != nil {
			return nil, err
		}
	}

	return &docker{
		client: c,
	}, nil
}

func (d *docker) Info(ctx context.Context) (system.Info, error) {
	return d.client.Info(ctx)
}

func (d *docker) Version(ctx context.Context) (types.Version, error) {
	return d.client.ServerVersion(ctx)
}

func (d *docker) Container() container.Container {
	return container.New(d.client)
}

func (d *docker) Image() image.Image {
	return image.New(d.client)
}

func (d *docker) Network() network.Network {
	return network.New(d.client)
}

func (d *docker) Volume() volume.Volume {
	return volume.New(d.client)
}
