package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make(map[*websocket.Conn]bool)

// WebSocket handler to broadcast pacing updates
func speechPacer(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	clients[conn] = true

	// Keep the connection open
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			delete(clients, conn)
			break
		}
	}
}

// Serves the frontend (index.html)
func serveUI(c *gin.Context) {
	c.File("index.html")
}

func main() {
	r := gin.Default()
	r.GET("/", serveUI)
	r.GET("/ws", speechPacer)

	log.Println("Server running on port 80...")
	log.Fatal(r.Run(":80"))
}
