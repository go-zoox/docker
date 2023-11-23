package image

import (
	"context"
	"os"

	"github.com/docker/cli/cli/streams"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/jsonmessage"
)

// BuildOption is the configuration for building a image
type BuildOption = types.ImageBuildOptions

// Build builds a image
func (i *image) Build(ctx context.Context, src string, opts ...func(opt *BuildOption)) error {
	opt := &BuildOption{}
	for _, o := range opts {
		o(opt)
	}

	tar, err := archive.TarWithOptions(src, &archive.TarOptions{})
	if err != nil {
		return err
	}

	response, err := i.client.ImageBuild(ctx, tar, *opt)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// io.Copy(os.Stdout, response.Body)
	if err := jsonmessage.DisplayJSONMessagesToStream(response.Body, streams.NewOut(os.Stdout), nil); err != nil {
		return err
	}

	return nil
}
