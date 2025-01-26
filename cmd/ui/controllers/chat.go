package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatController interface {
	ChatPage(c *gin.Context, html HTML)
}

type chatController struct{}

func NewChatController() ChatController {
	return &chatController{}
}

func (cc *chatController) ChatPage(c *gin.Context, html HTML) {
	if html == nil {
		html = c.HTML
	}
	html(http.StatusOK, "chat", gin.H{
		"title": "Chat Page",
		"ws":    "ws://localhost:3000/chat/ws",
	})
}
