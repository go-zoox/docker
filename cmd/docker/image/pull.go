package image

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/docker"
)

func Pull() *cli.Command {
	return &cli.Command{
		Name:  "pull",
		Usage: "Download an image from a registry",
		Flags: []cli.Flag{},
		Action: func(ctx *cli.Context) error {
			image := ctx.Args().First()
			if image == "" {
				return fmt.Errorf("image name is required")
			}

			client, err := docker.New()
			if err != nil {
				return err
			}

			return client.Image().Pull(ctx.Context, image)
		},
	}
}
