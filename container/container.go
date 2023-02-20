package container

import "github.com/docker/docker/api/types"

// Container is the docker container client interface
type Container interface {
	Run(cfg *RunConfig) error
	Logs(id string) ([]byte, error)
	Start(id string) error
	Stop(id string) error
	Remove(id string) error
	Inspect(id string) (*types.ContainerJSON, error)
	Status(id string) (*types.ContainerState, error)
}

type container struct {
}

// New creates a docker container client
func New() Container {
	return &container{}
}

func (c *container) Run(cfg *RunConfig) error {
	return Run(cfg)
}

func (c *container) Logs(id string) ([]byte, error) {
	return Logs(id)
}

func (c *container) Start(id string) error {
	return Start(id)
}

func (c *container) Stop(id string) error {
	return Stop(id)
}

func (c *container) Remove(id string) error {
	return Remove(id)
}

func (c *container) Inspect(id string) (*types.ContainerJSON, error) {
	return Inspect(id)
}

func (c *container) Status(id string) (*types.ContainerState, error) {
	return Status(id)
}
