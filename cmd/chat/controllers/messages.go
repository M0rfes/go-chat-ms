package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Message struct {
	From    string `json:"from"`
	Message string `json:"message"`
}

type MessagesController interface {
	HandleWebSocket(*gin.Context)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type messagesController struct{}

func NewMessagesController() MessagesController {
	return &messagesController{}
}

func (m *messagesController) HandleWebSocket(c *gin.Context) {
	userId, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userIdStr, ok := userId.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade WebSocket:", err)
		return
	}
	defer conn.Close()

	log.Println("Client connected:", conn.RemoteAddr())
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		message := &Message{
			From:    userIdStr,
			Message: string(p),
		}
		messageStr, err := json.Marshal(message)
		if err != nil {
			log.Println("Marshal error:", err)
			break
		}
		if err := conn.WriteMessage(messageType, []byte(messageStr)); err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}
