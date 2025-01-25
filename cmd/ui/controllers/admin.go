package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin", gin.H{
		"title": "Admin Page",
	})
}
