package container

import (
	"fmt"
	"os"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/go-zoox/cli"
	"github.com/go-zoox/docker"
	co "github.com/go-zoox/docker/container"
)

func Logs() *cli.Command {
	return &cli.Command{
		Name:  "logs",
		Usage: "Fetch the logs of a container",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "details",
				Usage: "Show extra details provided to logs",
			},
			&cli.BoolFlag{
				Name:    "follow",
				Aliases: []string{"f"},
				Usage:   "Follow log output",
			},
			&cli.StringFlag{
				Name:  "since",
				Usage: "Show logs since timestamp (e.g. 2013-01-02T13:23:37) or relative (e.g. 42m for 42 minutes)",
			},
			&cli.StringFlag{
				Name:    "tail",
				Aliases: []string{"n"},
				Usage:   "Number of lines to show from the end of the logs",
				Value:   "all",
			},
			&cli.BoolFlag{
				Name:    "timestamps",
				Aliases: []string{"t"},
				Usage:   "Show timestamps",
			},
			&cli.StringFlag{
				Name:  "until",
				Usage: "Show logs before a timestamp (e.g. 2013-01-02T13:23:37) or relative (e.g. 42m for 42 minutes)",
			},
		},
		Action: func(ctx *cli.Context) error {
			client, err := docker.New()
			if err != nil {
				return err
			}

			containerID := ctx.Args().First()
			if containerID == "" {
				return fmt.Errorf("container id is required")
			}

			logs, err := client.Container().Logs(ctx.Context, containerID, func(opt *co.LogsConfig) {
				opt.Follow = ctx.Bool("follow")
				opt.Timestamps = ctx.Bool("timestamps")
				opt.Tail = ctx.String("tail")
				opt.Since = ctx.String("since")
				opt.Until = ctx.String("until")
				opt.Details = ctx.Bool("details")
			})
			if err != nil {
				return err
			}
			defer logs.Close()

			if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, logs); err != nil {
				return err
			}

			return nil
		},
	}
}
