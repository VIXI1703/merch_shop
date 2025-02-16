package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "merch_shop/internal/model"
	"merch_shop/internal/provider"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestJWTAuthMiddleware(t *testing.T) {
	secret := "test_secret"
	auth := provider.NewJWTAuth([]byte(secret), time.Hour)
	middleware := JWTAuthMiddleware(auth)

	t.Run("NoAuthorizationHeader", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)

		middleware(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Authorization header required")
	})

	t.Run("InvalidTokenFormat", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.Header.Set("Authorization", "InvalidToken")

		middleware(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Authorization header required")
	})

	t.Run("InvalidToken", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.Header.Set("Authorization", "Bearer invalid.token")

		middleware(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Token invalid")
	})

	t.Run("ValidToken", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)

		token, _ := auth.GenerateToken(123)
		c.Request.Header.Set("Authorization", "Bearer "+token)

		middleware(c)

		claims, exists := GetUser(c)
		assert.True(t, exists)
		assert.Equal(t, uint(123), claims.UserId)
		assert.False(t, c.IsAborted())
	})
}

func TestGetUser(t *testing.T) {
	t.Run("NoUserInContext", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		_, exists := GetUser(c)
		assert.False(t, exists)
	})

	t.Run("ValidUserInContext", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		expectedClaims := provider.UserClaims{UserId: 123}
		c.Set("user", expectedClaims)

		claims, exists := GetUser(c)
		assert.True(t, exists)
		assert.Equal(t, expectedClaims, claims)
	})
}
