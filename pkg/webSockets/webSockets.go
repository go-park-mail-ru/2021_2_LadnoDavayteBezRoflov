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
			fmt.Println("SHIT: ", err)
			return
		}
		defer func(conn *websocket.Conn) {
			err := conn.Close()
			if err != nil {
				fmt.Println("SHITx2: ", err)
			}
		}(conn)

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

		mux.Lock()
		err = conn.ReadJSON(&inputData)
		if err != nil {
			fmt.Println("HERE: ", err)
			return
		}
		mux.Unlock()

		for userID, userConnections := range connections {
			for _, connection := range userConnections {
				if userID == uid {
					continue
				}
				mux.Lock()
				err = connection.WriteJSON(&inputData)
				if err != nil {
					fmt.Println("HERE2: ", err)
					mux.Unlock()
					break
				}
				mux.Unlock()
			}
		}
	}
}
