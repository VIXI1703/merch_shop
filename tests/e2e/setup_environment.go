package e2e

import (
	"bytes"
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
