package network

import (
	"sort"

	"github.com/go-zoox/cli"
	"github.com/go-zoox/docker"
	"github.com/go-zoox/docker/cmd/cli/table"
	"github.com/go-zoox/docker/network"
)

func List() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Usage:   "List docker networks",
		Aliases: []string{"ls"},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "Show all networks (default shows just running)",
			},
		},
		Action: func(ctx *cli.Context) error {
			client, err := docker.New()
			if err != nil {
				return err
			}

			networks, err := client.Network().List(ctx.Context, func(opt *network.ListOption) {
				// opt.All = ctx.Bool("all")
			})
			if err != nil {
				return err
			}

			columns := []table.Column{
				{Key: "id", Title: "Network ID", Width: 12, DisableEllipsis: true},
				{Key: "name", Title: "Name", Width: 28},
				{Key: "driver", Title: "DRIVER", Width: 12},
				{Key: "scope", Title: "SCOPE", Width: 18},
				{Key: "gateway", Title: "GATEWAY", Width: 18},
				{Key: "subnet", Title: "SUBNET", Width: 18},
			}

			sort.Slice(networks, func(i, j int) bool {
				return networks[i].Created.After(networks[j].Created)
			})

			rows := []map[string]string{}
			for _, network := range networks {
				ipamGateway := "-"
				ipamSubnet := "-"
				if len(network.IPAM.Config) > 0 {
					ipamGateway = network.IPAM.Config[0].Gateway
					ipamSubnet = network.IPAM.Config[0].Subnet
				}

				rows = append(rows, map[string]string{
					"id":      network.ID,
					"name":    network.Name,
					"driver":  network.Driver,
					"scope":   network.Scope,
					"gateway": ipamGateway,
					"subnet":  ipamSubnet,
				})
			}

			return table.Table(columns, rows)
		},
	}
}
