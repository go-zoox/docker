package image

import (
	"strings"

	"github.com/go-zoox/cli"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/docker"
	"github.com/go-zoox/docker/image"
)

func Build() *cli.Command {
	return &cli.Command{
		Name:  "build",
		Usage: "Build an image from a Dockerfile",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file",
				Usage:   "Name of the Dockerfile (Default is 'PATH/Dockerfile')",
				Aliases: []string{"f"},
			},
			&cli.StringSliceFlag{
				Name:    "tag",
				Usage:   "Name and optionally a tag in the 'name:tag' format",
				Aliases: []string{"t"},
			},
			&cli.StringFlag{
				Name:  "builkd-arg",
				Usage: "Set build-time variables",
			},
			&cli.StringFlag{
				Name:  "platform",
				Usage: "Set platform if server is multi-platform capable",
			},
		},
		Action: func(ctx *cli.Context) error {
			src := ctx.Args().First()
			if src == "" {
				return fmt.Errorf("src path is required")
			}

			client, err := docker.New()
			if err != nil {
				return err
			}

			return client.Image().Build(ctx.Context, src, func(opt *image.BuildOption) {
				opt.Dockerfile = ctx.String("file")
				opt.Tags = ctx.StringSlice("tag")

				opt.BuildArgs = map[string]*string{}
				for _, buildArg := range ctx.StringSlice("build-arg") {
					if buildArg == "" {
						continue
					}

					buildArgKV := strings.SplitN(buildArg, "=", 2)
					key := buildArgKV[0]
					value := ""
					if len(buildArgKV) > 1 {
						value = buildArgKV[1]
					}
					opt.BuildArgs[key] = &value
				}

				opt.Platform = ctx.String("platform")
			})
		},
	}
}
