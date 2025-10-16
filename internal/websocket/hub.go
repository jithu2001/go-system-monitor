package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jithu2001/go-system-monitor/internal/monitor"
)

type Hub struct {
	mu         sync.Mutex
	clients    map[*websocket.Conn]bool
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn
	broadcast  chan []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]bool),
		Register:   make(chan *websocket.Conn),
		Unregister: make(chan *websocket.Conn),
		broadcast:  make(chan []byte),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.Register:
			h.mu.Lock()
			h.clients[conn] = true
			h.mu.Unlock()
			log.Println("Client Connected!!")

		case conn := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.clients[conn]; ok {
				delete(h.clients, conn)
				conn.Close()
				log.Println("Client Disconnected!!")
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.Lock()
			for conn := range h.clients {
				if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
					delete(h.clients, conn)
					conn.Close()
					log.Println("Client Disconnected!!")
				}
			}
			h.mu.Unlock()
		}
	}
}
func (h *Hub) BroadcastStats() {
	for {
		stats, err := monitor.CollectStats()
		if err != nil {
			log.Println("Error Collecting Stats:", err)
			time.Sleep(2 * time.Second)
			continue
		}
		data, _ := json.Marshal(stats)
		h.broadcast <- data
		time.Sleep(1 * time.Second)
	}
}
