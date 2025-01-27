package controllers

import "github.com/gin-gonic/gin"

type HealthController interface {
	Health(*gin.Context)
}

type healthController struct{}

func NewHealthController() HealthController {
	return &healthController{}
}

func (h *healthController) Health(c *gin.Context) {
	c.JSON(200, gin.H{"status": "chat up"})
}
