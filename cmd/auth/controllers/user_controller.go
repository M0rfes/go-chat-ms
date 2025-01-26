package controllers

import (
	"net/http"
	"time"

	"github.com/M0rfes/go-chat-ms/auth/services"
	"github.com/gin-gonic/gin"
)

type UserControllers interface {
	Login(*gin.Context)
	Refresh(*gin.Context)
}

type userControllers struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) UserControllers {
	return &userControllers{
		userService: userService,
	}
}

func (u *userControllers) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := u.userService.Login(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", resp.Token, 300*int(time.Second), "/", "localhost", false, true)
	c.SetCookie("refresh-token", resp.RefreshToken, 30*24*int(time.Hour), "/", "localhost", false, true)

	c.JSON(http.StatusOK, resp)
}

func (u *userControllers) Refresh(c *gin.Context) {
	refreshToken := c.Request.Header.Get("refresh-token")
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token is required"})
		return
	}

	resp, err := u.userService.Refresh(refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", resp.Token, 60*60*24, "/", "localhost", false, true)

	c.JSON(200, resp)
}
