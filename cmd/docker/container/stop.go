package container

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/docker"
)

func Stop() *cli.Command {
	return &cli.Command{
		Name:  "stop",
		Usage: "stop one or more running containers",
		Flags: []cli.Flag{},
		Action: func(ctx *cli.Context) error {
			client, err := docker.New()
			if err != nil {
				return err
			}

			containerID := ctx.Args().First()
			if containerID == "" {
				return fmt.Errorf("conatiner id is required")
			}

			if err := client.Container().Stop(ctx.Context, containerID); err != nil {
				return err
			}

			fmt.Println(containerID)
			return nil
		},
	}
}
