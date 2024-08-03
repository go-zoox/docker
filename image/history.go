package image

import (
	"context"

	"github.com/go-zoox/docker/entity"
)

// History returns the changes in an image in history format.
func (i *image) History(ctx context.Context, id string) ([]entity.ImageHistory, error) {
	return i.client.ImageHistory(ctx, id)
}
