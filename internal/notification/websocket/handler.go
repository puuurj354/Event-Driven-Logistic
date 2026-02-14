package websocket

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	gorillaWs "github.com/gorilla/websocket"
)

var upgrader = gorillaWs.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWs(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("❌ WebSocket upgrade gagal: %v", err)
			return
		}

		client := &Client{
			Hub:  hub,
			Send: make(chan []byte, 256),
		}
		hub.register <- client

		go writePump(conn, client)

		readPump(conn, client)
	}
}

func readPump(conn *gorillaWs.Conn, client *Client) {
	defer func() {
		client.Hub.unregister <- client
		conn.Close()
	}()

	conn.SetReadLimit(512)
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func writePump(conn *gorillaWs.Conn, client *Client) {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {

				conn.WriteMessage(gorillaWs.CloseMessage, []byte{})
				return
			}
			err := conn.WriteMessage(gorillaWs.TextMessage, message)
			if err != nil {
				return
			}

		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteMessage(gorillaWs.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
