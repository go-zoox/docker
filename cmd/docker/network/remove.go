package network

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/docker"
)

func Remove() *cli.Command {
	return &cli.Command{
		Name:  "remove",
		Usage: "Remove a network",
		Action: func(ctx *cli.Context) error {
			networkID := ctx.Args().First()
			if networkID == "" {
				return fmt.Errorf("network id is required")
			}

			client, err := docker.New()
			if err != nil {
				return err
			}

			if err := client.Network().Remove(ctx.Context, networkID); err != nil {
				return err
			}

			fmt.Println(networkID)
			return nil
		},
	}
}
