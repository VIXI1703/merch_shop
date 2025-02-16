package main

import (
	"context"
	"fmt"
	"merch_shop/internal"
	"merch_shop/internal/config"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	app.Start(ctx, cfg)
}
