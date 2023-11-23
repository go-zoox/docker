package volume

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/docker"
	"github.com/go-zoox/docker/cmd/cli/table"
	"github.com/go-zoox/docker/volume"
)

func List() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Usage:   "List docker volumes",
		Aliases: []string{"ls"},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "Show all volumes (default shows just running)",
			},
		},
		Action: func(ctx *cli.Context) error {
			client, err := docker.New()
			if err != nil {
				return err
			}

			volumes, err := client.Volume().List(ctx.Context, func(opt *volume.ListOption) {
				// opt.All = ctx.Bool("all")
			})
			if err != nil {
				return err
			}

			columns := []table.Column{
				// {Key: "id", Title: "Volume ID", Width: 12, DisableEllipsis: true},
				{Key: "name", Title: "Name", Width: 64},
				{Key: "driver", Title: "DRIVER", Width: 12},
				{Key: "scope", Title: "SCOPE", Width: 18},
				{Key: "created_at", Title: "CREATED_AT", Width: 32},
			}

			// sort.Slice(volumes, func(i, j int) bool {
			// 	return volumes[i].CreatedAt.After(volumes[j].Created)
			// })

			rows := []map[string]string{}
			for _, volume := range volumes {

				rows = append(rows, map[string]string{
					// "id":      volume.ID,
					"name":       volume.Name,
					"driver":     volume.Driver,
					"scope":      volume.Scope,
					"created_at": volume.CreatedAt,
				})
			}

			return table.Table(columns, rows)
		},
	}
}
