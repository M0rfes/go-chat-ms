package services

import (
	"time"

	pkg "github.com/M0rfes/go-chat-ms/pkg/token"
)

type UserService interface {
	LoginCheck(token, refreshToken string) (*LoginCheckResponse, error)
}

type userService struct {
	tokenService pkg.Token
}

func NewUserService(tokenService pkg.Token) UserService {
	return &userService{
		tokenService: tokenService,
	}
}

type LoginCheckResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

func (u *userService) LoginCheck(token, refreshToken string) (*LoginCheckResponse, error) {
	claims, err := u.tokenService.Validate(token)
	if err != nil {
		// refresh token
		claims, err = u.tokenService.Validate(refreshToken)
		if err != nil {
			return nil, err
		}
		token, err := u.tokenService.Sign(&pkg.Claims{
			UserID: claims.UserID,
		}, 300*time.Second)
		if err != nil {
			return nil, err
		}
		refreshToken, err := u.tokenService.Sign(&pkg.Claims{
			UserID: claims.UserID,
		}, 30*24*time.Hour)
		if err != nil {
			return nil, err
		}
		return &LoginCheckResponse{
			Token:        token,
			RefreshToken: refreshToken,
		}, nil
	}
	return &LoginCheckResponse{
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}
