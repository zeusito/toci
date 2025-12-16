package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/zeusito/toci/internal/healthcheck/handlers"
	"github.com/zeusito/toci/pkg/config"
	"github.com/zeusito/toci/pkg/logger"
	"github.com/zeusito/toci/pkg/router"
)

func main() {
	// Parse flags
	cfgPath := flag.String("config", "resources/config.toml", "Path to the configuration file")
	flag.Parse()

	// Setup logger
	logger.MustConfigure()

	// Load config
	myConfig, err := config.LoadConfigurations(*cfgPath)
	if err != nil {
		log.Fatalf("Error loading configurations: %s", err.Error())
	}

	// Init router
	theRouter := router.NewHTTPRouter(myConfig.Server)

	// Health Controller
	_ = handlers.NewHealthController(theRouter.Mux)

	// Start server in background
	go theRouter.Start()

	// Graceful shutdown
	gracefulShutdown(theRouter)
}

func gracefulShutdown(myRouter *router.HTTPRouter) {
	// Wait for the interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	// Signal acquired, starting to shut down all systems
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	myRouter.Shutdown(ctx)
}
