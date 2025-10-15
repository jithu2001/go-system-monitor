package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type SystemInfo struct {
	CPUUsage    float64
	MemoryUsage float64
	DiskUsage   float64
}

func main() {
	statsChan := make(chan SystemInfo)

	go func() {
		for {
			stats := collectStats()
			statsChan <- stats
			time.Sleep(1 * time.Second)
		}
	}()
	for s := range statsChan {
		// Clear screen (capital J) and move cursor to home
		fmt.Printf("\033[2J\033[H")
		fmt.Println(" System Resource Monitor")
		fmt.Println("==========================")
		fmt.Printf("CPU Usage: %.2f%%\n", s.CPUUsage)
		fmt.Printf("Memory Usage: %.2f%%\n", s.MemoryUsage)
		fmt.Printf("Disk Usage: %.2f%%\n", s.DiskUsage)
	}
}

func collectStats() SystemInfo {
	cpuPercent, _ := cpu.Percent(0, false)
	memStats, _ := mem.VirtualMemory()
	diskStats, _ := disk.Usage("/")

	return SystemInfo{
		CPUUsage:    cpuPercent[0],
		MemoryUsage: memStats.UsedPercent,
		DiskUsage:   diskStats.UsedPercent,
	}
}
