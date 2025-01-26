package controllers

import (
	"github.com/M0rfes/go-chat-ms/auth/services"
	"github.com/gin-gonic/gin"
)

type AdminControllers interface {
	Login(c *gin.Context)
	Refresh(c *gin.Context)
}

type adminControllers struct {
	adminService services.AdminService
}

func NewAdminController(adminService services.AdminService) AdminControllers {
	return &adminControllers{
		adminService: adminService,
	}
}

func (a *adminControllers) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	res, err := a.adminService.Login(req)
	if err != nil {
		if _, ok := err.(*services.InvalidCredentialsError); ok {
			c.JSON(401, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", res.Token, 60*60*24, "/", "localhost", false, true)
	c.SetCookie("refresh-token", res.RefreshToken, 60*60*24*30, "/", "localhost", false, true)

	c.JSON(200, res)
}

func (a *adminControllers) Refresh(c *gin.Context) {
	refreshToken := c.GetHeader("Authorization")
	if refreshToken == "" {
		c.JSON(401, gin.H{"error": "missing refresh token"})
		return
	}

	res, err := a.adminService.Refresh(refreshToken)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", res.Token, 60*60*24, "/", "localhost", false, true)
	c.SetCookie("refresh-token", res.RefreshToken, 60*60*24*30, "/", "localhost", false, true)

	c.JSON(200, res)
}
