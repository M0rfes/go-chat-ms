package services

import (
	"time"

	pkg "github.com/M0rfes/go-chat-ms/pkg/token"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type UserService interface {
	Login(LoginRequest) (*LoginResponse, error)
	Refresh(string) (*LoginResponse, error)
}

type userService struct {
	tokenService pkg.Token
}

func NewUserService(tokenService pkg.Token) UserService {
	return &userService{
		tokenService: tokenService,
	}
}

func (s *userService) Login(req LoginRequest) (*LoginResponse, error) {
	claim := &pkg.Claims{
		UserID: req.Username,
	}
	token, err := s.tokenService.Sign(claim, 300*time.Second)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.tokenService.Sign(claim, 30*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (s *userService) Refresh(refreshToken string) (*LoginResponse, error) {
	claims, err := s.tokenService.Validate(refreshToken)
	if err != nil {
		return nil, err
	}

	token, err := s.tokenService.Sign(claims, 300*time.Second)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.tokenService.Sign(claims, 30*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:        token,
		RefreshToken: newRefreshToken,
	}, nil
}
