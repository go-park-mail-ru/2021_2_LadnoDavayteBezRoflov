package webSockets

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader    websocket.Upgrader
	connections map[uint]*websocket.Conn
)

func WebSocketsHandler(c *gin.Context) {
	var inputData struct {
		BID uint `json:"message"`
	}

	uid, exists := c.Get("uid")
	if exists {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		connections[uid.(uint)] = conn

		for {
			err := conn.ReadJSON(&inputData)
			if err != nil {
				break
			}

			for _, connection := range connections {
				err = connection.WriteJSON(&inputData)
				if err != nil {
					break
				}
			}
		}
	}
}
