package image

import (
	"github.com/docker/go-units"
	"github.com/go-zoox/cli"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/docker"
)

func Prune() *cli.Command {
	return &cli.Command{
		Name:  "prune",
		Usage: "Remove unused images",
		Flags: []cli.Flag{},
		Action: func(ctx *cli.Context) error {
			client, err := docker.New()
			if err != nil {
				return err
			}

			report, err := client.Image().Prune(ctx.Context)
			if err != nil {
				return err
			}

			lines := []string{}
			for _, image := range report.ImagesDeleted {
				if image.Deleted != "" {
					lines = append(lines, fmt.Sprintf("deleted: %s\n", image.Deleted))
				}

				if image.Untagged != "" {
					lines = append(lines, fmt.Sprintf("untagged: %s\n", image.Untagged))
				}
			}

			if len(lines) != 0 {
				lines = append(lines, "\n")
			}

			lines = append(lines, fmt.Sprintf("Total reclaimed space: %s\n", units.BytesSize(float64(report.SpaceReclaimed))))

			for _, line := range lines {
				fmt.Printf(line)
			}

			return nil
		},
	}
}
