package container

import (
	"os"
	"text/template"

	"github.com/go-zoox/cli"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/docker"
	co "github.com/go-zoox/docker/container"
)

func Inspect() *cli.Command {
	return &cli.Command{
		Name:  "inspect",
		Usage: "Inspect a container",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Usage:   "Format the output using the given go template",
				Aliases: []string{"f"},
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

			response, err := client.Container().Inspect(ctx.Context, containerID, func(opt *co.InspectOptions) {
				//
			})
			if err != nil {
				return err
			}

			if format := ctx.String("format"); format != "" {
				template, err := template.New("inspect").Parse(format)
				if err != nil {
					return err
				}

				return template.Execute(os.Stdout, response)
			}

			return fmt.PrintJSON(response)
		},
	}
}
