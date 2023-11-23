package volume

import (
	"context"

	vo "github.com/docker/docker/api/types/volume"
)

type CreateOption struct {
	// Name of the volume driver to use.
	Driver string `json:"Driver,omitempty"`

	// A mapping of driver options and values. These options are
	// passed directly to the driver and are driver specific.
	//
	DriverOpts map[string]string `json:"DriverOpts,omitempty"`

	// User-defined key/value metadata.
	Labels map[string]string `json:"Labels,omitempty"`
}

func (n *volume) Create(ctx context.Context, name string, opts ...func(*CreateOption)) (vo.Volume, error) {
	opt := &CreateOption{}
	for _, o := range opts {
		o(opt)
	}

	return n.client.VolumeCreate(ctx, vo.CreateOptions{
		Name:       name,
		Driver:     opt.Driver,
		DriverOpts: opt.DriverOpts,
		Labels:     opt.Labels,
	})
}
