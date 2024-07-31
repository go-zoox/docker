package image

import (
	"context"
	"os"

	"github.com/docker/cli/cli/streams"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/jsonmessage"
)

// BuildConfig is the configuration for building a image
type BuildConfig = types.ImageBuildOptions

// Build builds a image
func (i *image) Build(ctx context.Context, src string, opts ...func(cfg *BuildConfig)) error {
	cfg := &BuildConfig{}
	for _, o := range opts {
		o(cfg)
	}

	tar, err := archive.TarWithOptions(src, &archive.TarOptions{})
	if err != nil {
		return err
	}

	response, err := i.client.ImageBuild(ctx, tar, *cfg)
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
