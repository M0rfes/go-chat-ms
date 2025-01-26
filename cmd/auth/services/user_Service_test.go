package services_test

import (
	"testing"
	"time"

	"github.com/M0rfes/go-chat-ms/auth/services"
	pkg "github.com/M0rfes/go-chat-ms/pkg/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserTokenService struct {
	mock.Mock
}

func (m *mockUserTokenService) Sign(claims *pkg.Claims, duration time.Duration) (string, error) {
	args := m.Called(claims, duration)
	return args.String(0), args.Error(1)
}

func (m *mockUserTokenService) Validate(tokenString string) (*pkg.Claims, error) {
	args := m.Called(tokenString)
	return args.Get(0).(*pkg.Claims), args.Error(1)
}

func TestUserService_Login(t *testing.T) {
	t.Run("signs token with valid claims", func(t *testing.T) {
		mockToken := new(mockUserTokenService)
		userService := services.NewUserService(mockToken)

		mockToken.On("Sign", mock.Anything, 300*time.Second).Return("token", nil).Once()
		mockToken.On("Sign", mock.Anything, 30*24*time.Hour).Return("refreshToken", nil).Once()

		req := services.LoginRequest{
			Username: "test",
			Password: "test",
		}
		resp, err := userService.Login(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "token", resp.Token)
		assert.Equal(t, "refreshToken", resp.RefreshToken)
		mockToken.AssertExpectations(t)
	})

	t.Run("returns error when token signing fails", func(t *testing.T) {
		mockToken := new(mockUserTokenService)
		userService := services.NewUserService(mockToken)

		mockToken.On("Sign", mock.Anything, 300*time.Second).Return("", assert.AnError).Once()

		req := services.LoginRequest{
			Username: "test",
			Password: "test",
		}
		resp, err := userService.Login(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		mockToken.AssertExpectations(t)
	})
}

func TestUserService_Refresh(t *testing.T) {
	t.Run("successfully refreshes tokens", func(t *testing.T) {
		mockToken := new(mockUserTokenService)
		userService := services.NewUserService(mockToken)

		// Mock the validation of refresh token
		mockToken.On("Validate", "oldRefreshToken").Return(&pkg.Claims{
			UserID: "test-user",
		}, nil).Once()

		// Mock the signing of new tokens
		mockToken.On("Sign", mock.MatchedBy(func(claims *pkg.Claims) bool {
			return claims.UserID == "test-user"
		}), 300*time.Second).Return("newToken", nil).Once()

		mockToken.On("Sign", mock.MatchedBy(func(claims *pkg.Claims) bool {
			return claims.UserID == "test-user"
		}), 30*24*time.Hour).Return("newRefreshToken", nil).Once()

		resp, err := userService.Refresh("oldRefreshToken")

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "newToken", resp.Token)
		assert.Equal(t, "newRefreshToken", resp.RefreshToken)
		mockToken.AssertExpectations(t)
	})

	t.Run("returns error when refresh token is invalid", func(t *testing.T) {
		mockToken := new(mockUserTokenService)
		userService := services.NewUserService(mockToken)

		mockToken.On("Validate", "invalidToken").Return(&pkg.Claims{}, assert.AnError).Once()

		resp, err := userService.Refresh("invalidToken")

		assert.Error(t, err)
		assert.Nil(t, resp)
		mockToken.AssertExpectations(t)
	})

	t.Run("returns error when signing new token fails", func(t *testing.T) {
		mockToken := new(mockUserTokenService)
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
