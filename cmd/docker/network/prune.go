package network

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/docker"
)

func Prune() *cli.Command {
	return &cli.Command{
		Name:  "prune",
		Usage: "Remove unused networks",
		Flags: []cli.Flag{},
		Action: func(ctx *cli.Context) error {
			client, err := docker.New()
			if err != nil {
				return err
			}

			report, err := client.Network().Prune(ctx.Context)
			if err != nil {
				return err
			}

			lines := []string{}
			for _, network := range report.NetworksDeleted {
				lines = append(lines, fmt.Sprintf("deleted: %s\n", network))
			}

			for _, line := range lines {
				fmt.Printf(line)
			}

			return nil
		},
	}
}
