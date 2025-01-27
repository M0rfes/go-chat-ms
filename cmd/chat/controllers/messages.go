package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	model "github.com/M0rfes/go-chat-ms/pkg/message-queue/models"
	producer "github.com/M0rfes/go-chat-ms/pkg/message-queue/producers"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type MessagesController interface {
	HandleWebSocket(*gin.Context)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type messagesController struct {
	kafka producer.Message
}

func NewMessagesController(kafka producer.Message) MessagesController {
	return &messagesController{
		kafka: kafka,
	}
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
		message := &model.Message{
			From:    userIdStr,
			Message: string(p),
		}
		if err := m.kafka.Publish(*message); err != nil {
			log.Println("Publish error:", err)
			break
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
