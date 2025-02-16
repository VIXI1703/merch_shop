package e2e

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"merch_shop/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPurchaseItemScenario(t *testing.T) {
	srv := createTestServer(t)
	token := authenticateUser(t, srv, "test_buyer")

	// Test initial balance
	t.Run("CheckInitialBalance", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/info", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		srv.Gin.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var info model.InfoResponse
		err := json.Unmarshal(w.Body.Bytes(), &info)
		assert.NoError(t, err)
		assert.Equal(t, startBalance, info.Coins)
	})

	// Test item purchase
	t.Run("PurchaseItem", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/buy/t-shirt", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		srv.Gin.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify balance and inventory
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/api/info", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		srv.Gin.ServeHTTP(w, req)

		var info model.InfoResponse
		err := json.Unmarshal(w.Body.Bytes(), &info)
		assert.NoError(t, err)

		assert.Equal(t, startBalance-80, info.Coins)
		assert.Len(t, info.Inventory, 1)
		assert.Equal(t, "t-shirt", info.Inventory[0].Name)
		assert.Equal(t, 1, info.Inventory[0].Quantity)
	})
}
