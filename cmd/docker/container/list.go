package container

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-zoox/cli"
	"github.com/go-zoox/datetime"
	"github.com/go-zoox/docker"
	"github.com/go-zoox/docker/cmd/cli/table"
	"github.com/go-zoox/docker/container"
)

func List() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Usage:   "List docker containers",
		Aliases: []string{"ls"},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "Show all containers (default shows just running)",
			},
		},
		Action: func(ctx *cli.Context) error {
			client, err := docker.New()
			if err != nil {
				return err
			}

			containers, err := client.Container().List(ctx.Context, func(opt *container.ListOptions) {
				opt.All = ctx.Bool("all")
			})
			if err != nil {
				return err
			}

			columns := []table.Column{
				{Key: "id", Title: "CONTAINER ID", Width: 12, DisableEllipsis: true},
				{Key: "image", Title: "IMAGE", Width: 28},
				{Key: "command", Title: "COMMAND", Width: 22},
				{Key: "created_at", Title: "CREATED", Width: 20},
				{Key: "status", Title: "STATUS", Width: 32},
				{Key: "ports", Title: "PORTS", Width: 28},
				{Key: "names", Title: "NAMES", Width: 24},
			}

			rows := []map[string]string{}
			for _, container := range containers {
				names := []string{}
				for _, name := range container.Names {
					names = append(names, strings.TrimPrefix(name, "/"))
				}

				// createdAt := datetime.FromTime(time.UnixMilli(container.Created * 1e3)).Format("YYYY-MM-DD HH:mm:ss")
				createdAt := datetime.FromTime(time.UnixMilli(container.Created * 1e3)).Ago()
				ports := []string{}
				for _, port := range container.Ports {
					if port.PublicPort != 0 {
						ports = append(ports, fmt.Sprintf("%d->%d/%s", port.PrivatePort, port.PublicPort, port.Type))
					} else {
						ports = append(ports, fmt.Sprintf("%d/%s", port.PrivatePort, port.Type))
					}
				}

				rows = append(rows, map[string]string{
					"id":         container.ID,
					"image":      container.Image,
					"command":    container.Command,
					"created_at": createdAt,
					"status":     container.Status,
					"ports":      strings.Join(ports, ","),
					"names":      strings.Join(names, ","),
				})
			}

			return table.Table(columns, rows)
		},
	}
}
