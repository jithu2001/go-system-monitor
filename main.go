package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type SystemInfo struct {
	CPUUsage    float64
	MemoryUsage float64
	DiskUsage   float64
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	fmt.Println(" Websocket Server Started On ws://localhost:8080/ws")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		stats := collectStats()
		data, _ := json.Marshal(stats)
		err := conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Println("Write Error:", err)
			break
		}
		time.Sleep(1 * time.Second)
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
