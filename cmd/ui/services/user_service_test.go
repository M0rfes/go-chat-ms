package services

import (
	"errors"
	"testing"
	"time"

	pkg "github.com/M0rfes/go-chat-ms/pkg/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) Validate(token string) (*pkg.Claims, error) {
	args := m.Called(token)
	return args.Get(0).(*pkg.Claims), args.Error(1)
}

func (m *MockTokenService) Sign(claims *pkg.Claims, expiration time.Duration) (string, error) {
	args := m.Called(claims, expiration)
	return args.String(0), args.Error(1)
}

func TestLoginCheck(t *testing.T) {
	mockTokenService := new(MockTokenService)
	userService := NewUserService(mockTokenService)

	t.Run("valid token", func(t *testing.T) {
		mockClaims := &pkg.Claims{UserID: "123"}
		mockTokenService.On("Validate", "validToken").Return(mockClaims, nil)

		response, err := userService.LoginCheck("validToken", "validRefreshToken")
		assert.NoError(t, err)
		assert.Equal(t, "validToken", response.Token)
		assert.Equal(t, "validRefreshToken", response.RefreshToken)
		mockTokenService.AssertExpectations(t)
	})

	t.Run("invalid token, valid refresh token", func(t *testing.T) {
		mockClaims := &pkg.Claims{UserID: "123"}
		mockTokenService.On("Validate", "invalidToken").Return(&pkg.Claims{}, errors.New("invalid token"))
		mockTokenService.On("Validate", "validRefreshToken").Return(mockClaims, nil)

		// More specific mock for access token signing
		mockTokenService.On("Sign",
			mock.Anything,
			time.Second*300,
		).Return("newToken", nil).Once()

		// More specific mock for refresh token signing
		mockTokenService.On("Sign",
			mock.Anything,
			time.Hour*24*30,
		).Return("newRefreshToken", nil).Once()

		response, err := userService.LoginCheck("invalidToken", "validRefreshToken")
		assert.NoError(t, err)
		assert.Equal(t, "newToken", response.Token)
		assert.Equal(t, "newRefreshToken", response.RefreshToken)
		mockTokenService.AssertExpectations(t)
	})

	t.Run("invalid token and refresh token", func(t *testing.T) {
		mockTokenService.On("Validate", "invalidToken").Return(&pkg.Claims{}, errors.New("invalid token"))
		mockTokenService.On("Validate", "invalidRefreshToken").Return(&pkg.Claims{}, errors.New("invalid refresh token"))

		response, err := userService.LoginCheck("invalidToken", "invalidRefreshToken")
		assert.Error(t, err)
		assert.Nil(t, response)
		mockTokenService.AssertExpectations(t)
	})
}
