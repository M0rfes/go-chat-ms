package services_test

import (
	"testing"
	"time"

	"github.com/M0rfes/go-chat-ms/auth/services"
	pkg "github.com/M0rfes/go-chat-ms/pkg/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAdminTokenService struct {
	mock.Mock
}

func (m *mockAdminTokenService) Sign(claims *pkg.Claims, duration time.Duration) (string, error) {
	args := m.Called(claims, duration)
	return args.String(0), args.Error(1)
}

func (m *mockAdminTokenService) Validate(tokenString string) (*pkg.Claims, error) {
	args := m.Called(tokenString)
	return args.Get(0).(*pkg.Claims), args.Error(1)
}

func TestAdminService_Login(t *testing.T) {
	t.Run("signs token with valid claims", func(t *testing.T) {
		mockToken := new(mockAdminTokenService)
		adminService := services.NewAdminService(mockToken)

		mockToken.On("Sign", mock.Anything, 300*time.Second).Return("token", nil).Once()
		mockToken.On("Sign", mock.Anything, 30*24*time.Hour).Return("refreshToken", nil).Once()

		req := services.LoginRequest{
			Username: "admin",
			Password: "admin",
		}
		resp, err := adminService.Login(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "token", resp.Token)
		assert.Equal(t, "refreshToken", resp.RefreshToken)
		mockToken.AssertExpectations(t)
	})

	t.Run("returns error when credentials are invalid", func(t *testing.T) {
		mockToken := new(mockAdminTokenService)
		adminService := services.NewAdminService(mockToken)

		req := services.LoginRequest{
			Username: "test",
			Password: "test",
		}
		resp, err := adminService.Login(req)

		assert.Error(t, err)
		if err, ok := err.(*services.InvalidCredentialsError); !ok {
			t.Errorf("expected InvalidCredentialsError, got %v", err)
		}
		assert.EqualError(t, err, "invalid credentials")
		assert.Nil(t, resp)
	})

	t.Run("returns error when token signing fails", func(t *testing.T) {
		mockToken := new(mockAdminTokenService)
		adminService := services.NewAdminService(mockToken)

		mockToken.On("Sign", mock.Anything, 300*time.Second).Return("", assert.AnError).Once()

		req := services.LoginRequest{
			Username: "admin",
			Password: "admin",
		}
		resp, err := adminService.Login(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestAdminService_Refresh(t *testing.T) {
	t.Run("successfully refreshes tokens", func(t *testing.T) {
		mockToken := new(mockAdminTokenService)
		adminService := services.NewAdminService(mockToken)

		// Mock the validation of refresh token
		mockToken.On("Validate", "oldRefreshToken").Return(&pkg.Claims{
			UserID: "admin",
		}, nil).Once()

		// Mock the signing of new tokens
		mockToken.On("Sign", mock.MatchedBy(func(claims *pkg.Claims) bool {
			return claims.UserID == "admin"
		}), 300*time.Second).Return("newToken", nil).Once()

		mockToken.On("Sign", mock.MatchedBy(func(claims *pkg.Claims) bool {
			return claims.UserID == "admin"
		}), 30*24*time.Hour).Return("newRefreshToken", nil).Once()

		resp, err := adminService.Refresh("oldRefreshToken")

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "newToken", resp.Token)
		assert.Equal(t, "newRefreshToken", resp.RefreshToken)
		mockToken.AssertExpectations(t)
	})

	t.Run("returns error when refresh token is invalid", func(t *testing.T) {
		mockToken := new(mockAdminTokenService)
		userService := services.NewUserService(mockToken)

		mockToken.On("Validate", "invalidToken").Return(&pkg.Claims{}, assert.AnError).Once()

		resp, err := userService.Refresh("invalidToken")

		assert.Error(t, err)
		assert.Nil(t, resp)
		mockToken.AssertExpectations(t)
	})

	t.Run("returns error when signing new token fails", func(t *testing.T) {
		mockToken := new(mockAdminTokenService)
		userService := services.NewUserService(mockToken)

		mockToken.On("Validate", "validToken").Return(&pkg.Claims{
			UserID: "test-user",
		}, nil).Once()

		mockToken.On("Sign", mock.Anything, 300*time.Second).Return("", assert.AnError).Once()

		resp, err := userService.Refresh("validToken")

		assert.Error(t, err)
		assert.Nil(t, resp)
		mockToken.AssertExpectations(t)
	})
}
