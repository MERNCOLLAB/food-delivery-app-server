package realtime

import "github.com/gin-gonic/gin"

func DeliveryLocationWS(c *gin.Context) {
	orderId := c.Param("id")


	conn, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Error(err)
		return
	}

	ws := &WS{conn: conn}

	clientsMu.Lock()
	if clients[orderId] == nil {
		clients[orderId] = make(map[*WS]bool)
	}
	clients[orderId][ws] = true
	clientsMu.Unlock()

	defer func() {
		clientsMu.Lock()
		delete(clients[orderId], ws)
		clientsMu.Unlock()
	}()

	for {
		var msg map[string]interface{}
		if err := conn.ReadJSON(&msg); err != nil {
			break
		}
		broadcastLocation(orderId, msg)
	}
}
