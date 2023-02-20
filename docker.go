package docker

import (
	"github.com/go-zoox/docker/container"
	"github.com/go-zoox/docker/image"
)

// Docker is the client of docker
type Docker interface {
	Container() container.Container
	Image() image.Image
}

type docker struct {
}

// New creates a docker client
func New() Docker {
	return &docker{}
}

func (d *docker) Container() container.Container {
	return container.New()
}

func (d *docker) Image() image.Image {
	return image.New()
}
