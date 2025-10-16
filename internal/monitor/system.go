package monitor

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type SystemInfo struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	DiskUsage   float64 `json:"disk_usage"`
}

func CollectStats() (SystemInfo, error) {
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		return SystemInfo{}, err
	}
	memStats, err := mem.VirtualMemory()
	if err != nil {
		return SystemInfo{}, err
	}
	diskStats, err := disk.Usage("/")
	if err != nil {
		return SystemInfo{}, err
	}

	return SystemInfo{
		CPUUsage:    cpuPercent[0],
		MemoryUsage: memStats.UsedPercent,
		DiskUsage:   diskStats.UsedPercent,
	}, nil
}
