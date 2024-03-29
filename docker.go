package docker

import (
	"github.com/docker/docker/client"
	"github.com/go-zoox/docker/container"
	"github.com/go-zoox/docker/image"
	"github.com/go-zoox/docker/network"
	"github.com/go-zoox/docker/volume"
)

// Docker is the client of docker
type Docker interface {
	Container() container.Container
	Image() image.Image
	Network() network.Network
	Volume() volume.Volume
}

type docker struct {
	client *client.Client
}

// Options ...
type Options struct {
	Client *client.Client
}

// New creates a docker client.
func New(opts ...func(opt *Options)) (d Docker, err error) {
	opt := &Options{}
	for _, o := range opts {
		o(opt)
	}

	if opt.Client == nil {
		opt.Client, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			return nil, err
		}
	}

	return &docker{
		client: opt.Client,
	}, nil
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
