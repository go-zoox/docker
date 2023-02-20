package image

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

// BuildConfig is the configuration for building a image
type BuildConfig struct {
	Name       string
	Dockerfile string
	Context    string
}

// Build builds a image
func Build(cfg *BuildConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	tar, err := archive.TarWithOptions(cfg.Context, &archive.TarOptions{})
	if err != nil {
		return err
	}

	tags := []string{
		cfg.Name,
	}

	resp, err := cli.ImageBuild(ctx, tar, types.ImageBuildOptions{
		Dockerfile: cfg.Dockerfile,
		Tags:       tags,
		Remove:     true,
	})
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
