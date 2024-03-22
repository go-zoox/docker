package container

import (
	"context"
	"io"

	tc "github.com/docker/docker/api/types/container"
)

type LogsConfig struct {
	Follow     bool
	Timestamps bool
	Tail       string
	Since      string
	Until      string
	Details    bool
}

// Logs retrieves the logs of a container
func (c *container) Logs(ctx context.Context, id string, opts ...func(opt *LogsConfig)) (io.ReadCloser, error) {
	opt := &LogsConfig{}
	for _, o := range opts {
		o(opt)
	}

	return c.client.ContainerLogs(ctx, id, tc.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     opt.Follow,
		Timestamps: opt.Timestamps,
		Tail:       opt.Tail,
		Since:      opt.Since,
		Until:      opt.Until,
		Details:    opt.Details,
	})
}
