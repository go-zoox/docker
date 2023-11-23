package volume

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/docker"
)

func Remove() *cli.Command {
	return &cli.Command{
		Name:  "remove",
		Usage: "Remove a volume",
		Action: func(ctx *cli.Context) error {
			volumeID := ctx.Args().First()
			if volumeID == "" {
				return fmt.Errorf("volume id is required")
			}

			client, err := docker.New()
			if err != nil {
				return err
			}

			if err := client.Volume().Remove(ctx.Context, volumeID); err != nil {
				return err
			}

			fmt.Println(volumeID)
			return nil
		},
	}
}
