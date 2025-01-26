package auth

import (
	"time"

	pkg "github.com/M0rfes/go-chat-ms/pkg/token"
	"github.com/gin-gonic/gin"
)

type onUnauthorized func(*gin.Context)

type AuthMiddleware interface {
	TokenValidationMiddleware(onUnauthorized) gin.HandlerFunc
}

type authMiddleware struct {
	tokenService pkg.Token
}

func NewAuthMiddleware(tokenService pkg.Token) AuthMiddleware {
	return &authMiddleware{
		tokenService: tokenService,
	}
}

func (a *authMiddleware) TokenValidationMiddleware(cb onUnauthorized) gin.HandlerFunc {
	return func(c *gin.Context) {

		token, err := c.Cookie("token")
		if err != nil {
			cb(c)
			return
		}
		_, err = a.tokenService.Validate(token)
		if err != nil {
			refreshtoken, err := c.Cookie("refresh-token")
			if err != nil {
				cb(c)
				return
			}
			claims, err := a.tokenService.Validate(refreshtoken)
			if err != nil {
				cb(c)
				return
			}
			token, err := a.tokenService.Sign(&pkg.Claims{
				UserID: claims.UserID,
			}, 300*time.Second)
			if err != nil {
				cb(c)
				return
			}
			c.SetCookie("token", token, 300*int(time.Second), "/", "localhost", false, true)
			refreshtoken, err = a.tokenService.Sign(&pkg.Claims{
				UserID: claims.UserID,
			}, 30*24*time.Hour)
			if err != nil {
				cb(c)
				return
			}
			c.SetCookie("refresh-token", refreshtoken, 30*24*int(time.Hour), "/", "localhost", false, true)
		}
		c.Next()
	}
}
