package main

import (
	"example.com/chat/db"
	"example.com/chat/models"
	"example.com/chat/routes"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

type Client struct {
	conn    *websocket.Conn
	channel string
	send    chan []byte
}

var clients = make(map[*Client]bool)
var broadcast = make(chan []byte)
var mutex = &sync.Mutex{}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsEndpoint(c *gin.Context) {
	channel := c.Query("channel")
	username := c.Query("username")
	if channel == "" || username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Channel or username not specified"})
		return
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not open websocket connection"})
		return
	}

	client := &Client{conn: ws, channel: channel, send: make(chan []byte)}

	mutex.Lock()
	clients[client] = true
	mutex.Unlock()

	go handleMessages(client, username)
	go client.writeMessages()
}

func handleMessages(client *Client, username string) {
	defer func() {
		mutex.Lock()
		delete(clients, client)
		mutex.Unlock()
		client.conn.Close()
	}()

	for {
		_, msg, err := client.conn.ReadMessage()
		if err != nil {
			return
		}

		// Save message to database
		channelID, _ := models.GetChannelIDByName(client.channel)
		userID, _ := models.GetUserIDByUsername(username)
		message := models.Message{
			ChannelID: channelID,
			UserID:    userID,
			Content:   string(msg),
			Timestamp: time.Now(),
		}
		models.SaveMessage(message)

		broadcast <- msg
	}
}

func (c *Client) writeMessages() {
	for msg := range c.send {
		c.conn.WriteMessage(websocket.TextMessage, msg)
	}
}

func handleBroadcast() {
	for {
		msg := <-broadcast
		mutex.Lock()
		for client := range clients {
			select {
			case client.send <- msg:
			default:
				close(client.send)
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}

func main() {
	// Initialize the database
	db.InitDB()

	// Create a new Gin router
	router := gin.Default()

	// Register the WebSocket endpoint
	router.GET("/ws", func(c *gin.Context) {
		wsEndpoint(c)
	})

	// Register other API routes
	routes.RegisterRoutes(router)

	// Start handling WebSocket broadcast
	go handleBroadcast()

	// Start the server on port 8080
	router.Run(":8080")
}
