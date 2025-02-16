package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"merch_shop/internal/provider"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Authenticate(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}

// Test Helpers

func createTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func setUserContext(c *gin.Context, userId uint) {
	claims := provider.UserClaims{UserId: userId}
	c.Set("user", claims)
}

func TestAuthHandler_Authenticate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		c, w := createTestContext()
		c.Request = httptest.NewRequest("POST", "/auth",
			strings.NewReader(`{"username":"alice","password":"secret"}`))

		mockService := new(MockUserService)
		mockService.On("Authenticate", "alice", "secret").Return("token123", nil)

		handler := NewAuthHandler(mockService)
		handler.Authenticate(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"token":"token123"`)
	})

	t.Run("InvalidRequest", func(t *testing.T) {
		c, w := createTestContext()
		c.Request = httptest.NewRequest("POST", "/auth",
			strings.NewReader(`{"invalid":"data"}`))

		handler := NewAuthHandler(new(MockUserService))
		handler.Authenticate(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Required fields")
	})

	t.Run("AuthFailure", func(t *testing.T) {
		c, w := createTestContext()
		c.Request = httptest.NewRequest("POST", "/auth",
			strings.NewReader(`{"username":"bob","password":"wrong"}`))

		mockService := new(MockUserService)
		mockService.On("Authenticate", "bob", "wrong").Return("", errors.New("invalid credentials"))

		handler := NewAuthHandler(mockService)
		handler.Authenticate(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid credentials")
	})
}
