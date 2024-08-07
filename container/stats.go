package container

import (
	"context"
	"encoding/json"
	"io"
)

type StatsOptions struct {
	Stream bool
}

// Status retrieves the status of a container.
func (c *container) Stats(ctx context.Context, id string, opts ...func(opt *StatsOptions)) (io.ReadCloser, error) {
	opt := &StatsOptions{}
	for _, o := range opts {
		o(opt)
	}

	stats, err := c.client.ContainerStats(ctx, id, opt.Stream)
	if err != nil {
		return nil, err
	}

	// return &stats, nil
	// return &StatsReadCloser{stats.Body}, nil
	return stats.Body, nil
}

type StatsReadCloser struct {
	io.ReadCloser
}

func (s *StatsReadCloser) Close() error {
	return s.ReadCloser.Close()
}

func (s *StatsReadCloser) Read(p []byte) (int, error) {
	buf := make([]byte, 32768)

	n, err := s.ReadCloser.Read(buf)
	if err != nil {
		return 0, err
	}

	data := &ContainerStats{}
	if err := json.Unmarshal(buf[:n], data); err != nil {
		return 0, err
	}

	networkRx, networkTx := data.Network()
	ioRead, ioWrite := data.BlockIO()
	stats := map[string]any{
		"cpu":            data.CPUPercentage(),
		"memory":         data.Memory(),
		"memory_percent": data.MemoryPercentage(),
		"network": map[string]any{
			"rx": networkRx,
			"tx": networkTx,
		},
		"blockio": map[string]any{
			"read":  ioRead,
			"write": ioWrite,
		},
		// "network": data.Network(),
		// "blockio": data.BlockIO(),
		// "pids":    data.PidsCurrent(),
	}
	if sb, err := json.Marshal(stats); err != nil {
		return 0, err
	} else {
		return copy(p, sb), nil
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
