package e2e

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"merch_shop/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCoinTransferScenario(t *testing.T) {
	srv := createTestServer(t)
	senderToken := authenticateUser(t, srv, "sender")
	receiverToken := authenticateUser(t, srv, "receiver")

	// Test coin transfer
	transferAmount := uint(200)
	sendReq := model.SendCoinRequest{
		ToUser: "receiver",
		Amount: transferAmount,
	}
	body, _ := json.Marshal(sendReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+senderToken)
	srv.Gin.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify balances
	t.Run("CheckSenderBalance", func(t *testing.T) {
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/api/info", nil)
		req.Header.Set("Authorization", "Bearer "+senderToken)
		srv.Gin.ServeHTTP(w, req)

		var info model.InfoResponse
		err := json.Unmarshal(w.Body.Bytes(), &info)
		assert.NoError(t, err)
		assert.Equal(t, startBalance-transferAmount, info.Coins)
	})

	t.Run("CheckReceiverBalance", func(t *testing.T) {
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/api/info", nil)
		req.Header.Set("Authorization", "Bearer "+receiverToken)
		srv.Gin.ServeHTTP(w, req)

		var info model.InfoResponse
		err := json.Unmarshal(w.Body.Bytes(), &info)
		assert.NoError(t, err)
		assert.Equal(t, startBalance+transferAmount, info.Coins)
	})

	// Verify transaction history
	t.Run("CheckTransactionHistory", func(t *testing.T) {
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/api/info", nil)
		req.Header.Set("Authorization", "Bearer "+senderToken)
		srv.Gin.ServeHTTP(w, req)

		var info model.InfoResponse
		err := json.Unmarshal(w.Body.Bytes(), &info)
		assert.NoError(t, err)

		assert.Len(t, info.CoinHistory.Sent, 1)
		assert.Equal(t, "receiver", info.CoinHistory.Sent[0].ToUser)
		assert.Equal(t, transferAmount, info.CoinHistory.Sent[0].Amount)
	})
}
