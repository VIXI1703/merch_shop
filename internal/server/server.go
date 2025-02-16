package server

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"merch_shop/internal/config"
	"merch_shop/internal/db"
	"net/http"
	"time"
)

type Server struct {
	Cfg *config.Config
	Gin *gin.Engine
	DB  *gorm.DB
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		Cfg: cfg,
		Gin: gin.Default(),
		DB:  db.SetupDB(&cfg.DB),
	}
}

func (server *Server) Run(ctx context.Context) {
	srv := &http.Server{
		Addr:    ":" + server.Cfg.HTTP.Port,
		Handler: server.Gin.Handler(),
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
