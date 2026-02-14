package main

import (
	"context"

	"github.com/tugascript/loan_calculator/calculation_api/internal/config"
	"github.com/tugascript/loan_calculator/calculation_api/internal/providers/logger"
	"github.com/tugascript/loan_calculator/calculation_api/internal/server"
)

func main() {
	log := logger.DefaultLogger()
	ctx := context.Background()

	log.InfoContext(ctx, "Loading configuration...")
	cfg := config.NewConfig(log, "./.env")

	log = logger.ConfigLogger(cfg.Logger(), cfg.ENV())

	log.InfoContext(ctx, "Building server...")
	server := server.New(ctx, log, cfg)
	log.InfoContext(ctx, "Server built")

	done := make(chan bool, 1)

	log.InfoContext(ctx, "Starting server...")
	go server.Start(ctx, cfg.Port())
	log.InfoContext(ctx, "Server started")

	log.InfoContext(ctx, "Closing server...")
	go server.Close(ctx, done)
	<-done
	log.InfoContext(ctx, "Server stopped")
}
