package volume

import (
	"os"
	"text/template"

	"github.com/go-zoox/cli"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/docker"
)

func Inspect() *cli.Command {
	return &cli.Command{
		Name:  "inspect",
		Usage: "Inspect a volume",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Usage:   "Format the output using the given go template",
				Aliases: []string{"f"},
			},
		},
		Action: func(ctx *cli.Context) error {
			volumeID := ctx.Args().First()
			if volumeID == "" {
				return fmt.Errorf("volume id is required")
			}

			client, err := docker.New()
			if err != nil {
				return err
			}

			response, err := client.Volume().Inspect(ctx.Context, volumeID)
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
