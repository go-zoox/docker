package volume

import (
	"strings"

	"github.com/go-zoox/cli"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/docker"
	"github.com/go-zoox/docker/volume"
)

func Create() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "Create a volume",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "driver",
				Usage:   "Driver to manage the Volume",
				Aliases: []string{"d"},
				Value:   "bridge",
			},
			&cli.StringSliceFlag{
				Name:  "driver-opt",
				Usage: "Set driver specific options",
			},
			&cli.StringFlag{
				Name:  "label",
				Usage: "Set metadata on a volume",
			},
		},
		Action: func(ctx *cli.Context) error {
			name := ctx.Args().First()
			if name == "" {
				return fmt.Errorf("volume name is required")
			}

			client, err := docker.New()
			if err != nil {
				return err
			}

			response, err := client.Volume().Create(ctx.Context, func(opt *volume.CreateOption) {
				opt.Name = name
				opt.Driver = ctx.String("driver")
				opt.DriverOpts = map[string]string{}
				for _, driverOpt := range ctx.StringSlice("driver-opt") {
					key, value := splitKeyValue(driverOpt)
					opt.DriverOpts[key] = value
				}

				opt.Labels = map[string]string{}
				for _, label := range ctx.StringSlice("label") {
					key, value := splitKeyValue(label)
					opt.Labels[key] = value
				}
			})
			if err != nil {
				return err
			}

			fmt.Println(response.Name)
			return nil
		},
	}
}

func splitKeyValue(label string) (key, value string) {
	kv := strings.SplitN(label, "=", 2)
	if len(kv) == 2 {
		key = kv[0]
		value = kv[1]
	} else {
		key = kv[0]
	}
	return
}
