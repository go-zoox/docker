package image

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/docker"
	"github.com/go-zoox/docker/image"
)

func Remove() *cli.Command {
	return &cli.Command{
		Name:  "remove",
		Usage: "Remove one or more images",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "force",
				Usage:   "Force removal of the image",
				Aliases: []string{"f"},
			},
			&cli.BoolFlag{
				Name:  "no-prune",
				Usage: "Do not delete untagged parents",
			},
		},
		Action: func(ctx *cli.Context) error {
			client, err := docker.New()
			if err != nil {
				return err
			}

			imageID := ctx.Args().First()
			if imageID == "" {
				return fmt.Errorf("container id is required")
			}

			imagesDeleted, err := client.Image().Remove(ctx.Context, imageID, func(opt *image.RemoveOption) {
				opt.Force = ctx.Bool("force")
				opt.PruneChildren = ctx.Bool("no-prune")
			})
			if err != nil {
				return err
			}

			lines := []string{}
			for _, image := range imagesDeleted {
				if image.Deleted != "" {
					lines = append(lines, fmt.Sprintf("deleted: %s\n", image.Deleted))
				}

				if image.Untagged != "" {
					lines = append(lines, fmt.Sprintf("untagged: %s\n", image.Untagged))
				}
			}

			for _, line := range lines {
				fmt.Printf(line)
			}

			return nil
		},
	}
}
