package volume

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/docker"
)

func Prune() *cli.Command {
	return &cli.Command{
		Name:  "prune",
		Usage: "Remove unused volumes",
		Flags: []cli.Flag{},
		Action: func(ctx *cli.Context) error {
			client, err := docker.New()
			if err != nil {
				return err
			}

			report, err := client.Volume().Prune(ctx.Context)
			if err != nil {
				return err
			}

			lines := []string{}
			for _, volume := range report.VolumesDeleted {
				lines = append(lines, fmt.Sprintf("deleted: %s\n", volume))
			}

			for _, line := range lines {
				fmt.Printf(line)
			}

			return nil
		},
	}
}
