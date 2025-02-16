package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"merch_shop/internal/config"
	"merch_shop/internal/db"
	"merch_shop/internal/model"
	"merch_shop/internal/server"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	startBalance  = uint(1000)
	testJWtSecret = "test_secret"
)

var testConfig = &config.Config{
	HTTP: config.HTTP{
		Port: "8080",
	},
	JWT: config.JWT{
		SigningKey: testJWtSecret,
		Duration:   24 * time.Hour,
	},
}

func createTestServer(t *testing.T) *server.Server {
	// Create in-memory database
	gormDB, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{SkipDefaultTransaction: true})
	assert.NoError(t, err)

	// Initialize database schema and seed data
	gormDB = db.InitDB(gormDB)

	// Create server with test config
	srv := &server.Server{
		Cfg: testConfig,
		Gin: gin.Default(),
		DB:  gormDB,
	}

	srv.ConfigureRoutes()

	return srv
}

func authenticateUser(t *testing.T, srv *server.Server, username string) string {
	authReq := model.AuthRequest{
		Username: username,
		Password: "password",
	}
	body, _ := json.Marshal(authReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(body))
	srv.Gin.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Authentication failed: %s", w.Body.String())

	var authResp model.AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &authResp)
	assert.NoError(t, err)

	return authResp.Token
}

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

func TestServerLifecycle(t *testing.T) {
	srv := createTestServer(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		srv.Run(ctx)
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Test server health
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/info", nil)
	srv.Gin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
