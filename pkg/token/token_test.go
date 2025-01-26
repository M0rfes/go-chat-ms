package token

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewAuth(t *testing.T) {

	t.Run("panics when secret not set", func(t *testing.T) {
		os.Unsetenv("TOKEN_SECRET_KEY")
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic but got none")
			}
		}()

		// Call the NewAuth function
		_ = NewAuth()
	})

	t.Run("signs token when secret is set", func(t *testing.T) {
		os.Setenv("TOKEN_SECRET_KEY", "test-secret")
		authService := NewAuth()
		assert.NotNil(t, authService)
	})
}
func TestSign(t *testing.T) {
	// Set up the environment variable for the secret key
	os.Setenv("TOKEN_SECRET_KEY", "test-secret")
	authService := NewAuth()

	t.Run("signs token with valid claims", func(t *testing.T) {
		claims := &Claims{
			UserID: "123",
		}
		token, err := authService.Sign(claims, 300*time.Second)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}

func TestValidate(t *testing.T) {
	// Set up the environment variable for the secret key
	os.Setenv("TOKEN_SECRET_KEY", "test-secret")
	authService := NewAuth()

	t.Run("validates token with correct secret", func(t *testing.T) {
		claims := &Claims{
			UserID: "123",
		}
		token, err := authService.Sign(claims, 300*time.Second)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		validatedClaims, err := authService.Validate(token)
		assert.NoError(t, err)
		assert.Equal(t, claims.UserID, validatedClaims.UserID)
	})

	t.Run("fails validation with wrong secret", func(t *testing.T) {
		// Change the secret key
		os.Setenv("TOKEN_SECRET_KEY", "test-secret")
		authService := NewAuth()

		claims := &Claims{
			UserID: "123",
		}

		// Sign the token with the correct secret
		token, err := authService.Sign(claims, 300*time.Second)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		os.Setenv("TOKEN_SECRET_KEY", "wrong-secret")
		authService = NewAuth()

		// Validate the token with the wrong secret
		_, err = authService.Validate(token)
		assert.Error(t, err)
	})

	t.Run("fails validation with expired token", func(t *testing.T) {
		os.Setenv("TOKEN_SECRET_KEY", "test-secret")
		authService := NewAuth()

		claims := &Claims{
			UserID: "123",
		}
		token, err := authService.Sign(claims, -1*time.Second)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		_, err = authService.Validate(token)
		assert.Error(t, err)
	})
}
