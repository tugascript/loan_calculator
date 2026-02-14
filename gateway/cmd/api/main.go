package main

import (
	"context"

	"github.com/tugascript/loan_calculator/gateway/internal/config"
	"github.com/tugascript/loan_calculator/gateway/internal/providers/logger"
	"github.com/tugascript/loan_calculator/gateway/internal/server"
)

func main() {
	log := logger.DefaultLogger()
	ctx := context.Background()

	log.InfoContext(ctx, "Loading configuration...")
	cfg := config.NewConfig(log, "./.env")
	log = logger.ConfigLogger(cfg.Logger(), cfg.ENV())
	log.InfoContext(ctx, "Configuration loaded")

	log.InfoContext(ctx, "Building server...")
	server := server.New(ctx, log, cfg)
	log.InfoContext(ctx, "Server built")

	done := make(chan bool)

	log.InfoContext(ctx, "Starting server...")
	go server.Start(ctx)
	log.InfoContext(ctx, "Server started")

	log.InfoContext(ctx, "Closing server...")
	go server.Close(ctx, done)
	<-done
	log.InfoContext(ctx, "Server closed")
}
