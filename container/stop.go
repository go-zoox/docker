package container

import (
	"context"

	dc "github.com/docker/docker/api/types/container"
)

// StopOptions holds the options to stop or restart a container.
type StopOptions struct {
	// Timeout (optional) is the timeout (in seconds) to wait for the container
	// to stop gracefully before forcibly terminating it with SIGKILL.
	//
	// - Use nil to use the default timeout (10 seconds).
	// - Use '-1' to wait indefinitely.
	// - Use '0' to not wait for the container to exit gracefully, and
	//   immediately proceeds to forcibly terminating the container.
	// - Other positive values are used as timeout (in seconds).
	Timeout int `json:",omitempty"`

	// Signal (optional) is the signal to send to the container to (gracefully)
	// stop it before forcibly terminating the container with SIGKILL after the
	// timeout expires. If not value is set, the default (SIGTERM) is used.
	Signal string `json:"signal,omitempty"`
}

// Stop stops a container
func (c *container) Stop(ctx context.Context, id string, opts ...func(opt *StopOptions)) error {
	opt := &StopOptions{
		Timeout: 10,
	}
	for _, o := range opts {
		o(opt)
	}

	return c.client.ContainerStop(ctx, id, dc.StopOptions{
		Timeout: &opt.Timeout,
		Signal:  opt.Signal,
	})
}
