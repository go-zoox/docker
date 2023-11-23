package image

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/docker"
)

func Tag() *cli.Command {
	return &cli.Command{
		Name:  "tag",
		Usage: "Create a tag TARGET_IMAGE that refers to SOURCE_IMAGE",
		Action: func(ctx *cli.Context) error {
			source := ctx.Args().Get(0)
			if source == "" {
				return fmt.Errorf("source is required")
			}

			target := ctx.Args().Get(1)
			if target == "" {
				return fmt.Errorf("target is required")
			}

			client, err := docker.New()
			if err != nil {
				return err
			}

			return client.Image().Tag(ctx.Context, source, target)
		},
	}
}
