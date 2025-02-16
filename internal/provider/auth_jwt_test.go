package provider

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWTAuth(t *testing.T) {
	secret := "test_secret"
	auth := NewJWTAuth([]byte(secret), time.Hour*24)

	t.Run("GenerateAndVerifyValidToken", func(t *testing.T) {
		userId := uint(123)
		token, err := auth.GenerateToken(userId)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		claims, err := auth.VerifyToken(token)
		assert.NoError(t, err)
		assert.Equal(t, userId, claims.UserId)
		assert.WithinDuration(t, time.Now().Add(time.Hour*24), claims.ExpiresAt.Time, time.Minute)
	})

	t.Run("VerifyInvalidToken", func(t *testing.T) {
		_, err := auth.VerifyToken("invalid.token.string")
		assert.Error(t, err)
	})

	t.Run("VerifyExpiredToken", func(t *testing.T) {
		expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjEyMywiRXhwaXJlc0F0IjoiMjAyMy0wMS0wMVQwMDowMDowMFoifQ.2g8m_CCAPgDmVtCF6gX_7h0D6z4q5J6w9V0y7vW7ZkQ"
		_, err := auth.VerifyToken(expiredToken)
		assert.Error(t, err)
	})

	t.Run("VerifyWrongSignature", func(t *testing.T) {
		wrongSecretToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjEyMywiRXhwaXJlc0F0IjoiMjAyMy0wMS0wMVQwMDowMDowMFoifQ.7B3NvqV7JXZ7t7D7Q7W7ZkQ"
		_, err := auth.VerifyToken(wrongSecretToken)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "signature")
	})
}
