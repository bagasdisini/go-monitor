package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"go-monitor/pkg/models"
	"log"
	"net/http"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan models.WebsiteStatus)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	clients[ws] = true

	for {
		var msg models.WebsiteStatus
		err := ws.ReadJSON(&msg)
		if err != nil {
			delete(clients, ws)
			break
		}
	}
}

func BroadcastStatus(status models.WebsiteStatus) {
	broadcast <- status
}

func handleMessages() {
	for {
		status := <-broadcast
		message, err := json.Marshal(status)
		if err != nil {
			continue
		}
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func init() {
	go handleMessages()
}
