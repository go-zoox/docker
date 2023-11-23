package volume

import (
	"context"
)

type RemoveOption struct {
	Force bool `json:"force"`
}

func (n *volume) Remove(ctx context.Context, id string, opts ...func(opt *RemoveOption)) error {
	opt := &RemoveOption{}
	for _, o := range opts {
		o(opt)
	}

	return n.client.VolumeRemove(ctx, id, opt.Force)
}
