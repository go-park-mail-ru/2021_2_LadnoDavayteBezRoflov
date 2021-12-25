package webSockets

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader    websocket.Upgrader
	connections map[uint][]*websocket.Conn
	mux         *sync.Mutex
)

func SetupWebSocketHandler() {
	connections = make(map[uint][]*websocket.Conn)
	mux = &sync.Mutex{}
	upgrader = websocket.Upgrader{
		HandshakeTimeout: 10 * time.Microsecond,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

func WebSocketsHandler(c *gin.Context) {
	var inputData struct {
		BID uint `json:"message"`
	}

	_uid, exists := c.Get("uid")
	uid := _uid.(uint)
	if exists {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer conn.Close()

		{
			mux.Lock()
			defer mux.Unlock()
			for _, elem := range connections[uid] {
				if elem == conn {
					fmt.Println("CONFLICT")
					return
				}
			}
			connections[uid] = append(connections[uid], conn)
		}

		defer func() {
			mux.Lock()
			defer mux.Unlock()
			for idx, elem := range connections[uid] {
				if elem == conn {
					connections[uid][idx] = connections[uid][len(connections[uid])-1]
					connections[uid] = connections[uid][:len(connections[uid])-1]
					return
				}
			}
		}()

		err = conn.ReadJSON(&inputData)
		if err != nil {
			fmt.Println(err)
			return
		}

		for userID, userConnections := range connections {
			mux.Lock()
			for _, connection := range userConnections {
				if userID == uid {
					continue
				}
				err = connection.WriteJSON(&inputData)
				if err != nil {
					fmt.Println(err)
					mux.Unlock()
					break
				}
			}
			mux.Unlock()
		}
	}
}
