package realtime

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WS struct {
	conn *websocket.Conn
}

var Upgrader = websocket.Upgrader{
	WriteBufferSize: 1024,
	ReadBufferSize:  1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	clients   = make(map[string]map[*WS]bool)
	clientsMu sync.Mutex
)

func broadcastLocation(orderID string, msg map[string]interface{}) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for client := range clients[orderID] {
		client.conn.WriteJSON(msg)
	}
}
