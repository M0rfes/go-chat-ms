package services

import (
	"time"

	pkg "github.com/M0rfes/go-chat-ms/pkg/token"
)

type AdminService interface {
	Login(LoginRequest) (*LoginResponse, error)
	Refresh(string) (*LoginResponse, error)
}

type adminService struct {
	tokenService pkg.Token
}

type InvalidCredentialsError struct{}

func (e *InvalidCredentialsError) Error() string {
	return "invalid credentials"
}

func NewAdminService(tokenService pkg.Token) AdminService {
	return &adminService{
		tokenService: tokenService,
	}
}

func (s *adminService) Login(req LoginRequest) (*LoginResponse, error) {
	if req.Username != "admin" || req.Password != "admin" {
		return nil, &InvalidCredentialsError{}
	}
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

func (s *adminService) Refresh(refreshToken string) (*LoginResponse, error) {
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
