package image

import (
	"strings"

	"github.com/docker/go-units"
	"github.com/go-zoox/cli"
	"github.com/go-zoox/datetime"
	"github.com/go-zoox/docker"
	"github.com/go-zoox/docker/cmd/cli/table"
	"github.com/go-zoox/docker/image"
)

func List() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Usage:   "List docker images",
		Aliases: []string{"ls"},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "Show all images (default shows just running)",
			},
		},
		Action: func(ctx *cli.Context) error {
			client, err := docker.New()
			if err != nil {
				return err
			}

			images, err := client.Image().List(ctx.Context, func(opt *image.ListOption) {
				opt.All = ctx.Bool("all")
			})
			if err != nil {
				return err
			}

			columns := []table.Column{
				{Key: "repository", Title: "REPOSITORY", Width: 53},
				{Key: "tag", Title: "TAG", Width: 18},
				{Key: "id", Title: "Image ID", Width: 12, DisableEllipsis: true},
				{Key: "created_at", Title: "CREATED", Width: 18},
				{Key: "size", Title: "Size", Width: 32},
			}

			rows := []map[string]string{}
			for _, image := range images {
				rows = append(rows, map[string]string{
					"repository": image.Name,
					"tag":        strings.Join(image.Tags, ","),
					"id":         image.ID,
					"created_at": datetime.FromTime(image.CreatedAt).Ago(),
					"size":       units.BytesSize(float64(image.Size)),
				})
			}

			return table.Table(columns, rows)
		},
	}
}
