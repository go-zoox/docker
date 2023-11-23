package main

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/docker/cmd/docker/container"
	"github.com/go-zoox/docker/cmd/docker/image"
	"github.com/go-zoox/docker/cmd/docker/network"
	"github.com/go-zoox/docker/cmd/docker/volume"
)

func main() {
	app := cli.NewMultipleProgram(&cli.MultipleProgramConfig{
		Name:  "docker",
		Usage: "cli of docker client",
	})

	app.Register("container", &cli.Command{
		Name:  "container",
		Usage: "Manage Containers",
		Subcommands: []*cli.Command{
			container.List(),
			container.Logs(),
			container.Inspect(),
			container.Stats(),
			container.Exec(),
			//
			container.Create(),
			container.Remove(),
			//
			container.Start(),
			container.Stop(),
			container.Restart(),
			//
			container.Run(),
		},
	})

	app.Register("image", &cli.Command{
		Name:  "image",
		Usage: "Manage Images",
		Subcommands: []*cli.Command{
			image.List(),
			image.Inspect(),
			//
			image.Remove(),
			//
			image.Pull(),
			//
			image.Build(),
			//
			image.Prune(),
			//
			image.Tag(),
		},
	})

	app.Register("network", &cli.Command{
		Name:  "network",
		Usage: "Manage Networks",
		Subcommands: []*cli.Command{
			network.List(),
			network.Inspect(),
			//
			network.Create(),
			network.Remove(),
			//
			network.Prune(),
		},
	})

	app.Register("volume", &cli.Command{
		Name:  "volume",
		Usage: "Manage Volumes",
		Subcommands: []*cli.Command{
			volume.List(),
			volume.Inspect(),
			//
			volume.Create(),
			volume.Remove(),
			//
			volume.Prune(),
		},
	})

	app.Run()
}
