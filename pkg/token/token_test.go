package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {
	authService := NewAuth("test-secret")

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
	authService := NewAuth("test-secret")

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
		authService := NewAuth("test-secret")

		claims := &Claims{
			UserID: "123",
		}

		// Sign the token with the correct secret
		token, err := authService.Sign(claims, 300*time.Second)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		authService = NewAuth("wrong-secret")

		// Validate the token with the wrong secret
		_, err = authService.Validate(token)
		assert.Error(t, err)
	})

	t.Run("fails validation with expired token", func(t *testing.T) {
		authService := NewAuth("test-secret")

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
