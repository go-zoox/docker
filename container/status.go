package container

import (
	"github.com/docker/docker/api/types"
)

// Status returns the status of a container
func Status(id string) (*types.ContainerState, error) {
	inspect, err := Inspect(id)
	if err != nil {
		return nil, err
	}

	return inspect.State, nil
}
