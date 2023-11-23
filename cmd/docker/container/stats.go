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

			stats, err := client.Container().Stats(ctx.Context, containerID, func(opt *co.StatsOptions) {
				opt.Stream = true
				if ctx.Bool("no-stream") {
					opt.Stream = false
				}
			})
			if err != nil {
				return err
			}
			defer stats.Body.Close()

			out := io.WriterWrapFunc(func(b []byte) (n int, err error) {
				data := &ContainerStats{}
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
						"name":        data.Name,
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
			if _, err := osio.Copy(out, stats.Body); err != nil {
				return err
			}

			return nil
		},
	}
}

type ContainerStats struct {
	Read      string `json:"read"`
	Preread   string `json:"preread"`
	PidsStats struct {
		Current int `json:"current"`
		// Limit   int64 `json:"limit"`
	} `json:"pids_stats"`
	BlkioStats struct {
		IoServiceBytesRecursive []struct {
			Major int    `json:"major"`
			Minor int    `json:"minor"`
			Op    string `json:"op"`
			Value int    `json:"value"`
		} `json:"io_service_bytes_recursive"`
		IoServicedRecursive []struct {
			Major int    `json:"major"`
			Minor int    `json:"minor"`
			Op    string `json:"op"`
			Value int    `json:"value"`
		} `json:"io_serviced_recursive"`
		IoQueueRecursive []struct {
			Major int    `json:"major"`
			Minor int    `json:"minor"`
			Op    string `json:"op"`
			Value int    `json:"value"`
		} `json:"io_queue_recursive"`
		IoServiceTimeRecursive []struct {
			Major int    `json:"major"`
			Minor int    `json:"minor"`
			Op    string `json:"op"`
			Value int    `json:"value"`
		} `json:"io_service_time_recursive"`
		IoWaitTimeRecursive []struct {
			Major int    `json:"major"`
			Minor int    `json:"minor"`
			Op    string `json:"op"`
			Value int    `json:"value"`
		} `json:"io_wait_time_recursive"`
		IoMergedRecursive []struct {
			Major int    `json:"major"`
			Minor int    `json:"minor"`
			Op    string `json:"op"`
			Value int    `json:"value"`
		} `json:"io_merged_recursive"`
		IoTimeRecursive []struct {
			Major int    `json:"major"`
			Minor int    `json:"minor"`
			Op    string `json:"op"`
			Value int    `json:"value"`
		} `json:"io_time_recursive"`
		SectorsRecursive []struct {
			Major int    `json:"major"`
			Minor int    `json:"minor"`
			Op    string `json:"op"`
			Value int    `json:"value"`
		} `json:"sectors_recursive"`
	} `json:"blkio_stats"`
	NumProcs     int `json:"num_procs"`
	StorageStats struct {
	} `json:"storage_stats"`
	PrecpuStats struct {
		CPUUsage struct {
			TotalUsage        int   `json:"total_usage"`
			PercpuUsage       []int `json:"percpu_usage"`
			UsageInKernelmode int   `json:"usage_in_kernelmode"`
			UsageInUsermode   int   `json:"usage_in_usermode"`
		} `json:"cpu_usage"`
		SystemCPUUsage int `json:"system_cpu_usage"`
		OnlineCPUs     int `json:"online_cpus"`
		ThrottlingData struct {
			Periods          int `json:"periods"`
			ThrottledPeriods int `json:"throttled_periods"`
			ThrottledTime    int `json:"throttled_time"`
		} `json:"throttling_data"`
	} `json:"precpu_stats"`
	CPUStats struct {
		CPUUsage struct {
			TotalUsage        int   `json:"total_usage"`
			UsageInKernelmode int   `json:"usage_in_kernelmode"`
			UsageInUsermode   int   `json:"usage_in_usermode"`
			PercpuUsage       []int `json:"percpu_usage"`
		} `json:"cpu_usage"`
		SystemCPUUsage int `json:"system_cpu_usage"`
		OnlineCPUs     int `json:"online_cpus"`
		ThrottlingData struct {
			Periods          int `json:"periods"`
			ThrottledPeriods int `json:"throttled_periods"`
			ThrottledTime    int `json:"throttled_time"`
		} `json:"throttling_data"`
	} `json:"cpu_stats"`
	MemoryStats struct {
		Usage    int `json:"usage"`
		MaxUsage int `json:"max_usage"`
		Stats    struct {
			ActiveAnon            int `json:"active_anon"`
			ActiveFile            int `json:"active_file"`
			Anon                  int `json:"anon"`
			AnonThp               int `json:"anon_thp"`
			File                  int `json:"file"`
			FileDity              int `json:"file_dirty"`
			FileMapped            int `json:"file_mapped"`
			FileWriteback         int `json:"file_writeback"`
			InactiveAnon          int `json:"inactive_anon"`
			InactiveFile          int `json:"inactive_file"`
			KernelStack           int `json:"kernel_stack"`
			Pgfault               int `json:"pgfault"`
			Pglazyfree            int `json:"pglazyfree"`
			Pglazyfreed           int `json:"pglazyfreed"`
			Pgmajfault            int `json:"pgmajfault"`
			Pgrefill              int `json:"pgrefill"`
			Pgscan                int `json:"pgscan"`
			Pgsteal               int `json:"pgsteal"`
			Shmem                 int `json:"shmem"`
			Slab                  int `json:"slab"`
			SlabReclaimable       int `json:"slab_reclaimable"`
			SlabUnreclaimable     int `json:"slab_unreclaimable"`
			Sock                  int `json:"sock"`
			ThpCollapseAlloc      int `json:"thp_collapse_alloc"`
			ThpFaultAlloc         int `json:"thp_fault_alloc"`
			Unevictable           int `json:"unevictable"`
			WorkingsetActivate    int `json:"workingset_activate"`
			WorkingsetNodereclaim int `json:"workingset_nodereclaim"`
			WorkingsetRefault     int `json:"workingset_refault"`
			//
			TotalInactiveFile int `json:"total_inactive_file"`
		} `json:"stats"`
		Limit int `json:"limit"`
	} `json:"memory_stats"`
	Networks map[string]struct {
		RxBytes   int `json:"rx_bytes"`
		RxPackets int `json:"rx_packets"`
		RxErrors  int `json:"rx_errors"`
		RxDropped int `json:"rx_dropped"`
		TxBytes   int `json:"tx_bytes"`
		TxPackets int `json:"tx_packets"`
		TxErrors  int `json:"tx_errors"`
		TxDropped int `json:"tx_dropped"`
	} `json:"networks"`
	Name string `json:"name"`
	ID   string `json:"id"`
}

func (c *ContainerStats) CPUPercentage() float64 {
	previousCPU := c.PrecpuStats.CPUUsage.TotalUsage
	previousSystem := c.PrecpuStats.SystemCPUUsage
	v := c
	var (
		cpuPercent = 0.0
		// calculate the change for the cpu usage of the container in between readings
		cpuDelta = float64(v.CPUStats.CPUUsage.TotalUsage) - float64(previousCPU)
		// calculate the change for the entire system between readings
		systemDelta = float64(v.CPUStats.SystemCPUUsage) - float64(previousSystem)
		onlineCPUs  = float64(v.CPUStats.OnlineCPUs)
	)

	if onlineCPUs == 0.0 {
		onlineCPUs = float64(len(v.CPUStats.CPUUsage.PercpuUsage))
	}

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * onlineCPUs * 100.0
	}

	return cpuPercent
}

func (c *ContainerStats) Memory() float64 {
	// cgroup v1
	isCgroup1 := c.MemoryStats.Stats.TotalInactiveFile != 0
	if v := c.MemoryStats.Stats.TotalInactiveFile; isCgroup1 && v < c.MemoryStats.Usage {
		return float64(c.MemoryStats.Usage - v)
	}

	// cgroup v2
	if v := c.MemoryStats.Stats.InactiveFile; !isCgroup1 && v < c.MemoryStats.Usage {
		return float64(c.MemoryStats.Usage - v)
	}

	return float64(c.MemoryStats.Usage)
}

func (c *ContainerStats) MemoryPercentage() float64 {
	usedNoCache := c.Memory()
	limit := c.MemoryLimit()
	if limit != 0 {
		return usedNoCache / limit * 100.0
	}

	return 0
}

func (c *ContainerStats) MemoryLimit() float64 {
	return float64(c.MemoryStats.Limit)
}

func (c *ContainerStats) Network() (rx float64, tx float64) {
	for _, v := range c.Networks {
		rx += float64(v.RxBytes)
		tx += float64(v.TxBytes)
	}
	return
}

func (c *ContainerStats) BlockIO() (read float64, write float64) {
	for _, bioEntry := range c.BlkioStats.IoServiceBytesRecursive {
		if len(bioEntry.Op) == 0 {
			continue
		}

		switch bioEntry.Op[0] {
		case 'r', 'R':
			read += float64(bioEntry.Value)
		case 'w', 'W':
			write += float64(bioEntry.Value)
		}
	}
	return
}

func (c *ContainerStats) PidsCurrent() int {
	return c.PidsStats.Current
}
