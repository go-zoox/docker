package container

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/docker"
	"github.com/go-zoox/docker/container"
)

func Remove() *cli.Command {
	return &cli.Command{
		Name:  "remove",
		Usage: "remove one or more containers",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Force the removal of a running container (uses SIGKILL)",
			},
			&cli.BoolFlag{
				Name:    "link",
				Aliases: []string{"l"},
				Usage:   "Remove the specified link",
			},
			&cli.BoolFlag{
				Name:    "volumes",
				Aliases: []string{"v"},
				Usage:   "Remove the volumes associated with the container",
			},
		},
		Action: func(ctx *cli.Context) error {
			client, err := docker.New()
			if err != nil {
				return err
			}

			containerID := ctx.Args().First()
			if containerID == "" {
				return fmt.Errorf("conatiner id is required")
			}

			err = client.Container().Remove(ctx.Context, containerID, func(opt *container.RemoveOptions) {
				opt.Force = ctx.Bool("force")
				opt.RemoveLinks = ctx.Bool("link")
				opt.RemoveVolumes = ctx.Bool("volumes")
			})
			if err != nil {
				return err
			}

			fmt.Println(containerID)
			return nil
		},
	}
}
