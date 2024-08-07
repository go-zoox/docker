package container

import (
	"encoding/json"
	osio "io"

	"github.com/docker/go-units"
	"github.com/go-zoox/core-utils/fmt"

	"github.com/go-zoox/cli"
	"github.com/go-zoox/core-utils/io"
	"github.com/go-zoox/docker"
	"github.com/go-zoox/docker/cmd/cli/table"
	co "github.com/go-zoox/docker/container"
)

func Stats() *cli.Command {
	return &cli.Command{
		Name:  "stats",
		Usage: "Display a live stream of container(s) resource usage statistics",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "no-stream",
				Usage: "Disable streaming stats and only pull the first result",
			},
		},
		Action: func(ctx *cli.Context) error {
			client, err := docker.New()
			if err != nil {
				return err
			}

			containerID := ctx.Args().First()
			if containerID == "" {
				return fmt.Errorf("container id is required")
			}

			reader, err := client.Container().Stats(ctx.Context, containerID, func(opt *co.StatsOptions) {
				opt.Stream = true
				if ctx.Bool("no-stream") {
					opt.Stream = false
				}
			})
			if err != nil {
				return err
			}
			defer reader.Close()

			out := io.WriterWrapFunc(func(b []byte) (n int, err error) {
				data := &co.ContainerStats{}
				if err := json.Unmarshal(b, data); err != nil {
					return 0, err
				}

				columns := []table.Column{
					{Key: "id", Title: "CONTAINER ID", Width: 12, DisableEllipsis: true},
					{Key: "name", Title: "NAME", Width: 24},
					{Key: "cpu", Title: "CPU %", Width: 8},
					{Key: "mem", Title: "MEM USAGE / LIMIT", Width: 24},
					{Key: "mem_percent", Title: "MEM %", Width: 8},
					{Key: "net_io", Title: "NET I/O", Width: 20, DisableEllipsis: true},
					{Key: "block_io", Title: "BLOCK I/O", Width: 24, DisableEllipsis: true},
					{Key: "pids", Title: "PIDS", Width: 8, DisableEllipsis: true},
				}
				networkRx, networkTx := data.Network()
				blockRead, blockWrite := data.BlockIO()
				dataSource := []map[string]string{
					{
						"id":          data.ID,
						"name":        data.Name[1:],
						"cpu":         fmt.Sprintf("%.2f%%", data.CPUPercentage()),
						"mem":         fmt.Sprintf("%s / %s", units.BytesSize(data.Memory()), units.BytesSize(float64(data.MemoryStats.Limit))),
						"mem_percent": fmt.Sprintf("%.2f%%", data.MemoryPercentage()),
						"net_io":      fmt.Sprintf("%s / %s", units.BytesSize(networkRx), units.BytesSize(networkTx)),
						"block_io":    fmt.Sprintf("%s / %s", units.BytesSize(blockRead), units.BytesSize(blockWrite)),
						"pids":        fmt.Sprintf("%d", data.PidsCurrent()),
					},
				}

				return len(b), table.Table(columns, dataSource, table.WithClearOption())
			})
			if _, err := osio.Copy(out, reader); err != nil {
				return err
			}

			return nil
		},
	}
}
