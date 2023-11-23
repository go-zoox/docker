package container

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/docker"
	"github.com/go-zoox/docker/container"
)

func Exec() *cli.Command {
	return &cli.Command{
		Name:  "exec",
		Usage: "Executes a command inside a container",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "detach",
				Aliases: []string{"d"},
				Usage:   "detach from the container's process",
			},
			&cli.BoolFlag{
				Name:    "tty",
				Aliases: []string{"t"},
				Usage:   "allocate a pseudo-TTY",
			},
		},
		Action: func(ctx *cli.Context) error {
			containerID := ctx.Args().First()
			if containerID == "" {
				return fmt.Errorf("container id is required")
			}

			cmd := ctx.Args().Tail()
			if len(cmd) == 0 {
				return fmt.Errorf("command is required")
			}

			client, err := docker.New()
			if err != nil {
				return err
			}

			return client.Container().Exec(ctx.Context, containerID, func(opt *container.ExecOptions) {
				// opt.Detach = ctx.Bool("detach")
				opt.Tty = true
				opt.Cmd = cmd
			})
		},
	}
}
