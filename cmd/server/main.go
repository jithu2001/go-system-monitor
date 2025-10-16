package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jithu2001/go-system-monitor/internal/config"
	ws "github.com/jithu2001/go-system-monitor/internal/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	hub := ws.NewHub()
	go hub.Run()
	go hub.BroadcastStats()
	port := config.GetPort()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}
		hub.Register <- conn

		go func() {
			defer func() { hub.Unregister <- conn }()
			for {
				if _, _, err := conn.NextReader(); err != nil {
					break
				}
			}
		}()
	})
	fmt.Println(" Websocket Server Started On ws://localhost:8080/ws")
	log.Fatal(http.ListenAndServe(port, nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {

}
