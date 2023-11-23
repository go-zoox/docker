package image

import "context"

func (i *image) Tag(ctx context.Context, source, target string) error {
	return i.client.ImageTag(ctx, source, target)
}
