package container

import (
	"io"
	"os"

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
			&cli.BoolFlag{
				Name:    "interactive",
				Aliases: []string{"i"},
				Usage:   "Keep STDIN open even if not attached",
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

			stream, err := client.Container().Exec(ctx.Context, containerID, func(opt *container.ExecOptions) {
				// opt.Detach = ctx.Bool("detach")
				opt.Tty = true
				opt.Cmd = cmd
			})
			if err != nil {
				return err
			}
			defer stream.Close()

			if ctx.Bool("interactive") {
				go io.Copy(stream, os.Stdin)
			}

			if _, err := io.Copy(os.Stdout, stream); err != nil {
				return err
			}

			return nil
		},
	}
}
