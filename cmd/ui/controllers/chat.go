package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ChatPage(c *gin.Context) {
	c.HTML(http.StatusOK, "chat", gin.H{
		"title": "Login Page",
	})
}
