package container

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/docker"
)

func Start() *cli.Command {
	return &cli.Command{
		Name:  "start",
		Usage: "start one or more stopped containers",
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

			if err := client.Container().Start(ctx.Context, containerID); err != nil {
				return err
			}

			fmt.Println(containerID)
			return nil
		},
	}
}
