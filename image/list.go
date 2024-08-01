package image

import (
	"context"

	"github.com/docker/docker/api/types/filters"
	dimage "github.com/docker/docker/api/types/image"
	"github.com/go-zoox/docker/entity"
)

type ListConfig struct {
	// All controls whether all images in the graph are filtered, or just
	// the heads.
	All bool

	// Filters is a JSON-encoded set of filter arguments.
	Filters filters.Args

	// SharedSize indicates whether the shared size of images should be computed.
	SharedSize bool

	// ContainerCount indicates whether container count should be computed.
	ContainerCount bool
}

func (i *image) List(ctx context.Context, opts ...func(opt *ListConfig)) (images []entity.Image, err error) {
	opt := &ListConfig{}
	for _, o := range opts {
		o(opt)
	}

	return i.client.ImageList(ctx, dimage.ListOptions{})

	// imagesX, err := i.client.ImageList(ctx, dimage.ListOptions{})
	// if err != nil {
	// 	return nil, err
	// }

	// for _, image := range imagesX {
	// 	if len(image.RepoTags) == 0 {
	// 		continue
	// 	}

	// 	// if image.RepoTags[0] == "<none>:<none>" {
	// 	// 	continue
	// 	// }

	// 	for _, repoTag := range image.RepoTags {
	// 		name := ""
	// 		tags := []string{}
	// 		if repoTag != "" {
	// 			tagXs := strings.Split(repoTag, ":")
	// 			tagLength := len(tagXs)
	// 			name = strings.Join(tagXs[0:tagLength-1], ":")
	// 			tags = append(tags, tagXs[tagLength-1])
	// 		}

	// 		images = append(images, entity.Image{
	// 			ID:        image.ID[7:],
	// 			Name:      name,
	// 			Tags:      tags,
	// 			Size:      image.Size,
	// 			CreatedAt: time.UnixMilli(image.Created * 1000),
	// 		})
	// 	}
	// }

	// return
}
