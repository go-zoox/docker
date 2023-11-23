package volume

import (
	"context"

	"github.com/docker/docker/api/types/filters"
	"github.com/go-zoox/docker/entity"
)

type ListOption struct {
	Filters filters.Args
}

func (n *volume) List(ctx context.Context, opts ...func(*ListOption)) ([]entity.Volume, error) {
	opt := &ListOption{}
	for _, o := range opts {
		o(opt)
	}

	respone, err := n.client.VolumeList(ctx, opt.Filters)
	if err != nil {
		return nil, err
	}

	volumes := make([]entity.Volume, 0, len(respone.Volumes))
	for _, v := range respone.Volumes {
		volumes = append(volumes, entity.Volume{
			Name:       v.Name,
			Driver:     v.Driver,
			Mountpoint: v.Mountpoint,
			CreatedAt:  v.CreatedAt,
			Labels:     v.Labels,
			Scope:      v.Scope,
		})
	}

	return volumes, nil
}
