package app

import (
	"context"
	"merch_shop/internal/config"
	"merch_shop/internal/server"
)

func Start(context context.Context, cfg config.Config) {
	app := server.NewServer(&cfg)

	app.ConfigureRoutes()

	app.Run(context)
}
