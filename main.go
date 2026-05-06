package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/zthiagovalle/demo-golang/src/core/config"
)

// @title demo-golang
// @version 1.0
// @description Microserviço de demonstração de Clean Architecture com Echo, Postgres e SNS/SQS
func main() {
	cfg := mustLoadConfig()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app := newApp(ctx, cfg)
	defer app.close()

	app.setupConsumers(ctx)
	app.setupControllers()

	if err := app.start(ctx); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func mustLoadConfig() *config.Config {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	return cfg
}
