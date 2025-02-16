package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"merch_shop/internal/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockTransactionService struct {
	mock.Mock
}

func (m *MockTransactionService) GetInfo(userId uint) (model.InfoResponse, error) {
	args := m.Called(userId)
	return args.Get(0).(model.InfoResponse), args.Error(1)
}

func (m *MockTransactionService) SendCoin(userId uint, toUser string, amount uint) error {
	args := m.Called(userId, toUser, amount)
	return args.Error(0)
}

func (m *MockTransactionService) BuyItem(userId uint, name string) error {
	args := m.Called(userId, name)
	return args.Error(0)
}

func TestTransactionHandler_GetInfo(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		c, w := createTestContext()
		setUserContext(c, 1)

		mockService := new(MockTransactionService)
		mockService.On("GetInfo", uint(1)).Return(model.InfoResponse{
			Coins: 1000,
		}, nil)

		handler := NewTransactionHandler(mockService)
		handler.GetInfo(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"coins":1000`)
	})

	t.Run("ServiceError", func(t *testing.T) {
		c, w := createTestContext()
		setUserContext(c, 1)

		mockService := new(MockTransactionService)
		mockService.On("GetInfo", uint(1)).Return(model.InfoResponse{}, errors.New("database error"))

		handler := NewTransactionHandler(mockService)
		handler.GetInfo(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "database error")
	})
}

func TestTransactionHandler_SendCoin(t *testing.T) {
	t.Run("ValidRequest", func(t *testing.T) {
		c, w := createTestContext()
		setUserContext(c, 1)
		c.Request = httptest.NewRequest("POST", "/sendCoin",
			strings.NewReader(`{"toUser":"bob","amount":100}`))
		c.Request.Header.Set("Content-Type", "application/json")

		mockService := new(MockTransactionService)
		mockService.On("SendCoin", uint(1), "bob", uint(100)).Return(nil)

		handler := NewTransactionHandler(mockService)
		handler.SendCoin(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("InvalidRequest", func(t *testing.T) {
		c, w := createTestContext()
		c.Request = httptest.NewRequest("POST", "/sendCoin",
			strings.NewReader(`{"invalid":"data"}`))

		handler := NewTransactionHandler(new(MockTransactionService))
		handler.SendCoin(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Required fields")
	})

	t.Run("ServiceError", func(t *testing.T) {
		c, w := createTestContext()
		setUserContext(c, 1)
		c.Request = httptest.NewRequest("POST", "/sendCoin",
			strings.NewReader(`{"toUser":"bob","amount":100}`))

		mockService := new(MockTransactionService)
		mockService.On("SendCoin", uint(1), "bob", uint(100)).Return(errors.New("insufficient balance"))

		handler := NewTransactionHandler(mockService)
		handler.SendCoin(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "insufficient balance")
	})
}

func TestTransactionHandler_BuyItem(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		c, w := createTestContext()
		setUserContext(c, 1)
		c.Params = gin.Params{{Key: "item", Value: "sword"}}

		mockService := new(MockTransactionService)
		mockService.On("BuyItem", uint(1), "sword").Return(nil)

		handler := NewTransactionHandler(mockService)
		handler.BuyItem(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("MissingItem", func(t *testing.T) {
		c, w := createTestContext()
		setUserContext(c, 1)
		c.Params = gin.Params{}

		handler := NewTransactionHandler(new(MockTransactionService))
		handler.BuyItem(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		c, w := createTestContext()
		setUserContext(c, 1)
		c.Params = gin.Params{{Key: "item", Value: "shield"}}

		mockService := new(MockTransactionService)
		mockService.On("BuyItem", uint(1), "shield").Return(errors.New("item not found"))

		handler := NewTransactionHandler(mockService)
		handler.BuyItem(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "item not found")
	})
}
